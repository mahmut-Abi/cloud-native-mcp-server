package middleware

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultOIDCHTTPTimeout = 5 * time.Second
	defaultOIDCCacheTTL    = 10 * time.Minute
)

type oidcDiscoveryDocument struct {
	Issuer  string `json:"issuer"`
	JWKSURI string `json:"jwks_uri"`
}

type oidcJWKSDocument struct {
	Keys []oidcJWK `json:"keys"`
}

type oidcJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`

	// RSA
	N string `json:"n"`
	E string `json:"e"`

	// EC
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type oidcVerifier struct {
	discoveryURL string
	issuerURL    string

	expectedIssuer   string
	expectedAudience string

	httpClient *http.Client
	cacheTTL   time.Duration

	mu                 sync.RWMutex
	discovery          oidcDiscoveryDocument
	discoveryFetched   time.Time
	signingKeys        map[string]interface{}
	signingKeysFetched time.Time
}

func isOIDCConfigured(config AuthConfig) bool {
	return strings.TrimSpace(config.OIDCDiscoveryURL) != "" || strings.TrimSpace(config.OIDCIssuerURL) != ""
}

func newOIDCVerifier(config AuthConfig) (*oidcVerifier, error) {
	discoveryURL, issuerURL, err := buildOIDCDiscoveryURL(config.OIDCIssuerURL, config.OIDCDiscoveryURL)
	if err != nil {
		return nil, err
	}

	timeout := defaultOIDCHTTPTimeout
	if config.OIDCHTTPTimeoutSec > 0 {
		timeout = time.Duration(config.OIDCHTTPTimeoutSec) * time.Second
	}

	cacheTTL := defaultOIDCCacheTTL
	if config.OIDCJWKSCacheTTLSec > 0 {
		cacheTTL = time.Duration(config.OIDCJWKSCacheTTLSec) * time.Second
	}

	expectedAudience := strings.TrimSpace(config.OIDCAudience)
	if expectedAudience == "" {
		expectedAudience = strings.TrimSpace(config.OIDCClientID)
	}

	return &oidcVerifier{
		discoveryURL:     discoveryURL,
		issuerURL:        issuerURL,
		expectedIssuer:   strings.TrimSpace(config.OIDCIssuer),
		expectedAudience: expectedAudience,
		httpClient:       &http.Client{Timeout: timeout},
		cacheTTL:         cacheTTL,
		signingKeys:      make(map[string]interface{}),
	}, nil
}

func buildOIDCDiscoveryURL(issuerURL, discoveryURL string) (string, string, error) {
	issuerURL = strings.TrimSpace(issuerURL)
	discoveryURL = strings.TrimSpace(discoveryURL)

	if discoveryURL != "" {
		parsed, err := url.Parse(discoveryURL)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return "", "", fmt.Errorf("invalid OIDC discovery URL: %q", discoveryURL)
		}
		return discoveryURL, issuerURL, nil
	}

	if issuerURL == "" {
		return "", "", fmt.Errorf("OIDC issuer URL or discovery URL is required")
	}

	parsedIssuer, err := url.Parse(issuerURL)
	if err != nil || parsedIssuer.Scheme == "" || parsedIssuer.Host == "" {
		return "", "", fmt.Errorf("invalid OIDC issuer URL: %q", issuerURL)
	}

	normalizedIssuer := strings.TrimRight(issuerURL, "/")
	return normalizedIssuer + "/.well-known/openid-configuration", normalizedIssuer, nil
}

func (v *oidcVerifier) VerifyToken(rawToken string) error {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return fmt.Errorf("empty bearer token")
	}

	parserOptions := []jwt.ParserOption{
		jwt.WithValidMethods([]string{"RS256", "RS384", "RS512", "PS256", "PS384", "PS512", "ES256", "ES384", "ES512"}),
		jwt.WithExpirationRequired(),
	}

	expectedIssuer := v.expectedIssuer
	if expectedIssuer == "" {
		discovery, err := v.getDiscoveryDocument(false)
		if err != nil {
			return fmt.Errorf("failed to resolve OIDC discovery document: %w", err)
		}
		expectedIssuer = strings.TrimSpace(discovery.Issuer)
	}
	if expectedIssuer != "" {
		parserOptions = append(parserOptions, jwt.WithIssuer(expectedIssuer))
	}

	if v.expectedAudience != "" {
		parserOptions = append(parserOptions, jwt.WithAudience(v.expectedAudience))
	}

	parsedToken, err := jwt.Parse(rawToken, v.keyFunc, parserOptions...)
	if err != nil {
		return fmt.Errorf("OIDC token validation failed: %w", err)
	}
	if !parsedToken.Valid {
		return fmt.Errorf("OIDC token is not valid")
	}

	return nil
}

func (v *oidcVerifier) keyFunc(token *jwt.Token) (interface{}, error) {
	kidValue, ok := token.Header["kid"]
	if !ok {
		return nil, fmt.Errorf("OIDC token header missing kid")
	}

	kid, ok := kidValue.(string)
	if !ok || strings.TrimSpace(kid) == "" {
		return nil, fmt.Errorf("OIDC token kid header is invalid")
	}

	signingKeys, err := v.getSigningKeys(false)
	if err != nil {
		return nil, err
	}
	if key, found := signingKeys[kid]; found {
		return key, nil
	}

	// Key rotation can happen at any time, retry with forced refresh.
	signingKeys, err = v.getSigningKeys(true)
	if err != nil {
		return nil, err
	}
	if key, found := signingKeys[kid]; found {
		return key, nil
	}

	return nil, fmt.Errorf("OIDC signing key not found for kid %q", kid)
}

func (v *oidcVerifier) getDiscoveryDocument(forceRefresh bool) (oidcDiscoveryDocument, error) {
	v.mu.RLock()
	cached := v.discovery
	fetchedAt := v.discoveryFetched
	v.mu.RUnlock()

	if !forceRefresh && cached.JWKSURI != "" && time.Since(fetchedAt) < v.cacheTTL {
		return cached, nil
	}

	discovery, err := v.fetchDiscoveryDocument()
	if err != nil {
		// Fall back to cached document if available.
		if cached.JWKSURI != "" {
			return cached, nil
		}
		return oidcDiscoveryDocument{}, err
	}

	v.mu.Lock()
	v.discovery = discovery
	v.discoveryFetched = time.Now()
	v.mu.Unlock()

	return discovery, nil
}

func (v *oidcVerifier) fetchDiscoveryDocument() (oidcDiscoveryDocument, error) {
	req, err := http.NewRequest(http.MethodGet, v.discoveryURL, nil)
	if err != nil {
		return oidcDiscoveryDocument{}, fmt.Errorf("build OIDC discovery request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return oidcDiscoveryDocument{}, fmt.Errorf("request OIDC discovery document failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return oidcDiscoveryDocument{}, fmt.Errorf("OIDC discovery returned status %d", resp.StatusCode)
	}

	var discovery oidcDiscoveryDocument
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&discovery); err != nil {
		return oidcDiscoveryDocument{}, fmt.Errorf("decode OIDC discovery document failed: %w", err)
	}

	discovery.JWKSURI = strings.TrimSpace(discovery.JWKSURI)
	discovery.Issuer = strings.TrimSpace(discovery.Issuer)
	if discovery.Issuer == "" && v.issuerURL != "" {
		discovery.Issuer = v.issuerURL
	}
	if discovery.JWKSURI == "" {
		return oidcDiscoveryDocument{}, fmt.Errorf("OIDC discovery document missing jwks_uri")
	}

	return discovery, nil
}

func (v *oidcVerifier) getSigningKeys(forceRefresh bool) (map[string]interface{}, error) {
	v.mu.RLock()
	cached := v.signingKeys
	fetchedAt := v.signingKeysFetched
	v.mu.RUnlock()

	if !forceRefresh && len(cached) > 0 && time.Since(fetchedAt) < v.cacheTTL {
		return cached, nil
	}

	discovery, err := v.getDiscoveryDocument(forceRefresh)
	if err != nil {
		return nil, err
	}

	signingKeys, err := v.fetchSigningKeys(discovery.JWKSURI)
	if err != nil {
		// Fall back to cached keys if available.
		if len(cached) > 0 {
			return cached, nil
		}
		return nil, err
	}

	v.mu.Lock()
	v.signingKeys = signingKeys
	v.signingKeysFetched = time.Now()
	v.mu.Unlock()

	return signingKeys, nil
}

func (v *oidcVerifier) fetchSigningKeys(jwksURI string) (map[string]interface{}, error) {
	req, err := http.NewRequest(http.MethodGet, jwksURI, nil)
	if err != nil {
		return nil, fmt.Errorf("build OIDC JWKS request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request OIDC JWKS failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("OIDC JWKS returned status %d", resp.StatusCode)
	}

	var jwks oidcJWKSDocument
	if err := json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("decode OIDC JWKS failed: %w", err)
	}

	signingKeys := make(map[string]interface{}, len(jwks.Keys))
	for _, jwk := range jwks.Keys {
		if strings.TrimSpace(jwk.Kid) == "" {
			continue
		}
		if jwk.Use != "" && jwk.Use != "sig" {
			continue
		}

		key, err := jwk.publicKey()
		if err != nil {
			continue
		}
		signingKeys[jwk.Kid] = key
	}

	if len(signingKeys) == 0 {
		return nil, fmt.Errorf("OIDC JWKS did not contain any usable signing keys")
	}

	return signingKeys, nil
}

func (j oidcJWK) publicKey() (interface{}, error) {
	switch strings.ToUpper(strings.TrimSpace(j.Kty)) {
	case "RSA":
		return parseRSAPublicKey(j.N, j.E)
	case "EC":
		return parseECPublicKey(j.Crv, j.X, j.Y)
	default:
		return nil, fmt.Errorf("unsupported JWK key type %q", j.Kty)
	}
}

func parseRSAPublicKey(modulus, exponent string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(modulus))
	if err != nil {
		return nil, fmt.Errorf("decode RSA modulus failed: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(exponent))
	if err != nil {
		return nil, fmt.Errorf("decode RSA exponent failed: %w", err)
	}
	if len(nBytes) == 0 || len(eBytes) == 0 {
		return nil, fmt.Errorf("invalid RSA key parameters")
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)
	if !e.IsInt64() {
		return nil, fmt.Errorf("RSA exponent is too large")
	}

	publicExponent := int(e.Int64())
	if publicExponent <= 0 {
		return nil, fmt.Errorf("invalid RSA exponent")
	}

	return &rsa.PublicKey{
		N: n,
		E: publicExponent,
	}, nil
}

func parseECPublicKey(curveName, xCoord, yCoord string) (*ecdsa.PublicKey, error) {
	xBytes, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(xCoord))
	if err != nil {
		return nil, fmt.Errorf("decode EC x coordinate failed: %w", err)
	}
	yBytes, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(yCoord))
	if err != nil {
		return nil, fmt.Errorf("decode EC y coordinate failed: %w", err)
	}
	if len(xBytes) == 0 || len(yBytes) == 0 {
		return nil, fmt.Errorf("invalid EC key coordinates")
	}

	var curve elliptic.Curve
	switch strings.TrimSpace(curveName) {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("unsupported EC curve %q", curveName)
	}

	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)
	if !curve.IsOnCurve(x, y) {
		return nil, fmt.Errorf("EC public key coordinates are not on curve")
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

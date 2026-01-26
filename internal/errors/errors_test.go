package errors

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	err := New(ErrCodeInvalidParam, "test error")
	if err.Code() != ErrCodeInvalidParam {
		t.Errorf("expected %s, got %s", ErrCodeInvalidParam, err.Code())
	}
}

func TestWrapError(t *testing.T) {
	orig := errors.New("original error")
	err := Wrap(orig, ErrCodeInternal, "wrapped")
	if !errors.Is(err, orig) {
		t.Error("wrapped error should contain original")
	}
}

func TestWithHTTPStatus(t *testing.T) {
	err := New(ErrCodeNotFound, "not found").WithHTTPStatus(404)
	if err.HTTPStatus != 404 {
		t.Errorf("expected 404, got %d", err.HTTPStatus)
	}
}

func TestMissingParamError(t *testing.T) {
	err := MissingParamError("name")
	if err.HTTPStatus != 400 {
		t.Errorf("expected 400, got %d", err.HTTPStatus)
	}
}

func TestIs(t *testing.T) {
	err := New(ErrCodeTimeout, "timeout")
	if !Is(err, ErrCodeTimeout) {
		t.Error("Is should return true for matching code")
	}
}

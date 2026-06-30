package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// These values are intended to be injected at build time via -ldflags:
//
//	-X main.version=<version>
//	-X main.commit=<git sha>
//	-X main.buildTime=<rfc3339 timestamp>
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

// BuildInfo represents immutable build/runtime metadata of the server binary.
type BuildInfo struct {
	Version   string
	Commit    string
	ShortCommit string
	BuildTime string
	GoVersion string
	OS        string
	Arch      string
}

// resolveBuildInfo returns build metadata from ldflags with runtime fallback.
func resolveBuildInfo() BuildInfo {
	info := BuildInfo{
		Version:   normalizeBuildValue(version, "dev"),
		Commit:    normalizeBuildValue(commit, "unknown"),
		BuildTime: normalizeBuildValue(buildTime, "unknown"),
		GoVersion: strings.TrimPrefix(runtime.Version(), "go"),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	// Fallback to runtime build metadata if ldflags are missing.
	if bi, ok := debug.ReadBuildInfo(); ok {
		if info.Version == "dev" {
			mainVersion := strings.TrimSpace(bi.Main.Version)
			if mainVersion != "" && mainVersion != "(devel)" {
				info.Version = mainVersion
			}
		}
		if info.Commit == "unknown" {
			if revision := buildSettingValue(bi.Settings, "vcs.revision"); revision != "" {
				info.Commit = revision
			}
		}
		if info.BuildTime == "unknown" {
			if timestamp := buildSettingValue(bi.Settings, "vcs.time"); timestamp != "" {
				info.BuildTime = timestamp
			}
		}
	}

	if len(info.Commit) >= 7 {
		info.ShortCommit = info.Commit[:7]
	} else if info.Commit != "unknown" {
		info.ShortCommit = info.Commit
	} else {
		info.ShortCommit = "unknown"
	}

	return info
}

func normalizeBuildValue(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func buildSettingValue(settings []debug.BuildSetting, key string) string {
	for _, s := range settings {
		if s.Key == key {
			return strings.TrimSpace(s.Value)
		}
	}
	return ""
}

// printStartupBuildInfo writes a startup banner to stdout as the first output.
func printStartupBuildInfo(info BuildInfo) {
	// Parse build time for display
	buildTimeDisplay := info.BuildTime
	if t, err := time.Parse(time.RFC3339, info.BuildTime); err == nil {
		buildTimeDisplay = t.UTC().Format("2006-01-02 15:04:05 UTC")
	}

	banner := fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║  ☸  Cloud Native MCP Server
║
║  Version    %s
║  Commit     %s
║  Built      %s
║  Runtime    %s  %s/%s
╚══════════════════════════════════════════════════════════════╝`,
		info.Version,
		info.ShortCommit,
		buildTimeDisplay,
		info.GoVersion, info.OS, info.Arch,
	)

	_, _ = fmt.Fprintln(os.Stderr, banner)
}

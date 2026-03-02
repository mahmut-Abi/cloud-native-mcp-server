package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
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
	BuildTime string
	GoVersion string
}

// resolveBuildInfo returns build metadata from ldflags with runtime fallback.
func resolveBuildInfo() BuildInfo {
	info := BuildInfo{
		Version:   normalizeBuildValue(version, "dev"),
		Commit:    normalizeBuildValue(commit, "unknown"),
		BuildTime: normalizeBuildValue(buildTime, "unknown"),
		GoVersion: runtime.Version(),
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

// printStartupBuildInfo writes a first-line startup banner to stderr.
func printStartupBuildInfo(info BuildInfo) {
	_, _ = fmt.Fprintf(
		os.Stderr,
		"cloud-native-mcp-server version=%s commit=%s build_time=%s go=%s\n",
		info.Version,
		info.Commit,
		info.BuildTime,
		info.GoVersion,
	)
}

package main

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// runGitCommand runs a git command and returns the output
func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// setupGOMAXPROCS sets reasonable defaults for GOMAXPROCS
func setupGOMAXPROCS() {
	nproc := runtime.GOMAXPROCS(0)
	if nproc < 4 {
		nproc = 4
	}
	logrus.Debugf("Running with GOMAXPROCS=%d", runtime.GOMAXPROCS(nproc))
}

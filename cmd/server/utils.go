package main

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

// setupGOMAXPROCS sets reasonable defaults for GOMAXPROCS
func setupGOMAXPROCS() {
	nproc := runtime.GOMAXPROCS(0)
	if nproc < 4 {
		nproc = 4
	}
	logrus.Debugf("Running with GOMAXPROCS=%d", runtime.GOMAXPROCS(nproc))
}

//+build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var targets = "linux/amd64 darwin/amd64 windows/amd64"

var goexe = "go"

const (
	packageName = "github.com/jerrinss5/go-jiraintegration"
)

var (
	curDir    string
	binDir    string
	distDir   string
	hash      string
	buildDate string
	buildNum  string
	version   string
)

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")

	curDir, err := os.Getwd()
	if err != nil {
		curDir = "." //hack
	}
	binDir = curDir
	distDir = path.Join(curDir, "dist")
}

var gox = sh.RunCmd("gox")

// Runs go mod download and creates a distributable executable
func Dist() error {
	mg.Deps(checkGox)
	fmt.Printf("[+] Cross compiling for: %q\n", targets)
	err := sh.RunWith(flagEnv(),
		"gox",
		"-parallel=3",
		"-output",
		"$DISTDIR/{{.Dir}}-{{.OS}}-{{.Arch}}",
		"--osarch=$TARGETS",
		"cmd/jiraintegration/main.go",
	)
	if err != nil {
		return err
	}
	fmt.Printf("[+] Cross compiled binaries in: %s\n", distDir)
	return nil
}

func checkGox() error {
	_, err := exec.LookPath("gox")
	if err != nil {
		return sh.Run(goexe, "get", "-u", "github.com/mitchellh/gox")
	}
	return nil
}

// set up environment variables
func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	if version = os.Getenv("VERSION"); version == "" {
		version = "dev"
	}
	if buildNum = os.Getenv("BUILDNUM"); buildNum == "" {
		buildNum = "local"
	}

	return map[string]string{
		"PKG":         packageName,
		"GOBIN":       binDir,
		"GIT_SHA":     hash,
		"DATE":        time.Now().Format("01/02/06"),
		"VERSION":     version,
		"BUILDNUM":    buildNum,
		"DISTDIR":     distDir,
		"CGO_ENABLED": "1", //bug: when this is disabled, DNS gets wonky
		"TARGETS":     targets,
	}
}

//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	slash       = string(filepath.Separator)
	BinDir      = "bin" + slash
	CoverageDir = "cov" + slash
)

func init() {
	rootDir, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		BinDir = filepath.Join(".", BinDir)
		CoverageDir = filepath.Join(".", CoverageDir)
	} else {
		BinDir = filepath.Join(strings.TrimSpace(string(rootDir)), BinDir)
		CoverageDir = filepath.Join(strings.TrimSpace(string(rootDir)), CoverageDir)

	}
}

func Clean() error {
	return sh.Rm(BinDir)
}

func Build() error {
	mg.Deps(Clean)

	if err := os.Mkdir(BinDir, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to build bin directory %s: %w", BinDir, err)
	}

	return sh.RunV("go", "build", "-o", BinDir, "")

}

type Test mg.Namespace

func (Test) Clean() error {
	return sh.Run(CoverageDir)
}

func (Test) Unit() error {
	mg.Deps(Clean)

	if err := os.Mkdir(CoverageDir, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to build coverage directory %s: %w", CoverageDir, err)
	}

	coverFile := filepath.Join(CoverageDir, "coverage.txt")
	coverProfile := fmt.Sprintf("-coverprofile=%s", coverFile)

	return sh.RunV("go", "test", "-v", "-short", "-timeout", "60s", "-count=1", coverProfile, "./...")
}

func (Test) HtmlCov() error {
	mg.Deps(Clean)

	if err := os.Mkdir(CoverageDir, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to build coverage directory %s: %w", CoverageDir, err)
	}

	coverFile := filepath.Join(CoverageDir, "coverage.txt")
	coverProfile := fmt.Sprintf("-coverprofile=%s", coverFile)

	err := sh.RunV("go", "test", "-v", "-short", "-timeout", "60s", "-count=1", coverProfile, "./...")
	_ = sh.RunV("go", "tool", "cover", fmt.Sprintf("-html=%s", coverFile))
	return err
}

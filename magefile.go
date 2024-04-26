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
	slash = string(filepath.Separator)

	//
	ccValidatorDir = "cc-validator"
	BinDir         = "bin" + slash
	CoverageDir    = "cov" + slash
	coverFile      = "coverage.txt"

	staticSiteDir = "static-site"
	tfPlanFile    = "static-site.tfplan"
)

func init() {
	rootDir := "."
	gitRepoRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil {
		rootDir = strings.TrimSpace(string(gitRepoRoot))
	}

	BinDir = filepath.Join(rootDir, BinDir)
	CoverageDir = filepath.Join(rootDir, CoverageDir)
	coverFile = filepath.Join(CoverageDir, coverFile)
	tfPlanFile = filepath.Join(rootDir, tfPlanFile)
}

type StaticSite mg.Namespace

func (StaticSite) Init() error {
	os.Chdir(staticSiteDir)
	defer os.Chdir("..")

	return sh.RunV("terraform", "init")
}

func (StaticSite) Build() error {
	mg.SerialDeps(StaticSite.Init)

	os.Chdir(staticSiteDir)
	defer os.Chdir("..")

	return sh.RunV("terraform", "plan", "-out", tfPlanFile)
}

func (StaticSite) Deploy() error {
	mg.SerialDeps(StaticSite.Build)

	os.Chdir(staticSiteDir)
	defer os.Chdir("..")

	return sh.RunV("terraform", "apply", tfPlanFile)
}

type CcValidator mg.Namespace

func (CcValidator) Clean() error {
	_ = sh.Rm(BinDir)
	_ = sh.Rm(CoverageDir)

	return nil
}

func (CcValidator) Build() error {
	mg.Deps(CcValidator.Clean)

	if err := os.Mkdir(BinDir, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to build bin directory %s: %w", BinDir, err)
	}

	os.Chdir(ccValidatorDir)
	defer os.Chdir("..")

	return sh.RunV("go", "build", "-o", BinDir, "")

}

func (CcValidator) Test() error {
	mg.Deps(CcValidator.Build)

	if err := os.Mkdir(CoverageDir, 0700); err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to build coverage directory %s: %w", CoverageDir, err)
	}

	coverProfile := fmt.Sprintf("-coverprofile=%s", coverFile)

	os.Chdir(ccValidatorDir)
	defer os.Chdir("..")

	return sh.RunV("go", "test", "-v", "-short", "-timeout", "60s", "-count=1", coverProfile, "./...")
}

func (CcValidator) Coverage() error {
	mg.Deps(CcValidator.Test)

	_ = sh.RunV("go", "tool", "cover", fmt.Sprintf("-html=%s", coverFile))
	return nil
}

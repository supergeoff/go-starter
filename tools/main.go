//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func run(workDir string, name string, args ...string) error { // Added workDir parameter
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if workDir != "" { // Set working directory if provided
		cmd.Dir = workDir
	}

	return cmd.Run()
}

// getWorkspaceModules remains the same as the corrected version
func getWorkspaceModules() ([]string, error) {
	workFilePath := filepath.Join("..", "go.work")
	data, err := os.ReadFile(workFilePath)
	if err != nil {
		wd, _ := os.Getwd()
		return nil, fmt.Errorf("failed to read %s (CWD: %s): %w", workFilePath, wd, err)
	}
	wf, err := modfile.ParseWork(filepath.Base(workFilePath), data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", workFilePath, err)
	}
	var modules []string
	for _, use := range wf.Use {
		modules = append(modules, use.Path)
	}
	return modules, nil
}

func Lint() error {
	fmt.Println("Reading modules from go.work (from tools/magefile.go)...")
	modules, err := getWorkspaceModules()
	if err != nil {
		return fmt.Errorf("could not get workspace modules: %w", err)
	}

	if len(modules) == 0 {
		fmt.Println("No modules found in go.work. Nothing to lint.")
		return nil
	}

	fmt.Printf("Found modules to lint: %v\n", modules)
	fmt.Println("Linting Go modules in workspace (from tools/magefile.go)...")

	for _, modulePath := range modules {
		fmt.Printf("===> Linting module: %s\n", modulePath)
		cleanModulePath := filepath.Clean(modulePath)

		// This is the actual directory where the module resides, relative to workspace root
		// e.g., "apps/client"
		// We need to construct the path to this directory relative to tools/
		actualModuleDir := filepath.Join("..", cleanModulePath) // e.g., ../apps/client

		args := []string{
			"run",
			"./...", // Lint all packages within the CWD
		}

		fmt.Printf("Running: (cd %s && golangci-lint %v)\n", actualModuleDir, args)
		// Call run with the target directory
		err := run(actualModuleDir, "golangci-lint", args...) // MODIFIED call
		if err != nil {
			// The error from run already includes command details if cmd.Run() fails
			return fmt.Errorf("golangci-lint failed for module %s: %w", cleanModulePath, err)
		}
		fmt.Printf("<=== Finished linting module: %s\n", cleanModulePath)
	}

	fmt.Println("All modules linted successfully (by tools/magefile.go).")
	return nil
}

func ServeBackend() error {
	actualModuleDir := filepath.Join("..", "apps/server")
	return run(actualModuleDir, "air", "-c", ".air.toml")
}

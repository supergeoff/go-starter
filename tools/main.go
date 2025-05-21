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
		modulePath = filepath.Clean(modulePath)

		// This is the actual directory where the module resides, relative to workspace root
		// e.g., "apps/client"
		// We need to construct the path to this directory relative to tools/
		relModuleDir := filepath.Join("..", modulePath) // e.g., ../apps/client

		args := []string{
			"run",
			"./...", // Lint all packages within the CWD
		}

		fmt.Printf("Running: (cd %s && golangci-lint %v)\n", relModuleDir, args)
		// Call run with the target directory
		err := run(relModuleDir, "golangci-lint", args...)
		if err != nil {
			// The error from run already includes command details if cmd.Run() fails
			return fmt.Errorf("golangci-lint failed for module %s: %w", modulePath, err)
		}
		fmt.Printf("<=== Finished linting module: %s\n", modulePath)
	}

	fmt.Println("All modules linted successfully (by tools/magefile.go).")
	return nil
}

func Serve(dirpath string) error {
	relModuleDir := filepath.Join("..", filepath.Clean(dirpath))
	return run(relModuleDir, "go", "tool", "air", "-c", ".air.toml")
}

// Install downloads and installs the Tailwind CSS CLI tool into the tools directory if not already present.
func Install() error {
	// finalPath is relative to the tools directory (where the magefile is)
	finalPath := "tailwindcss"

	// Check if the tool already exists and is executable in the tools directory
	if info, err := os.Stat(finalPath); err == nil {
		if !info.IsDir() && (info.Mode().Perm()&0o111 != 0) { // Check if it's a file and executable
			fmt.Printf(
				"Tailwind CSS (%s) already installed and executable in tools/. Skipping installation.\n",
				finalPath,
			)
			return nil
		}
		// File exists but is not a valid executable (e.g., a directory, or not executable).
		// We'll proceed to remove it and reinstall.
		fmt.Printf(
			"Found %s in tools/, but it's not a valid executable file or has wrong permissions. Proceeding with re-installation.\n",
			finalPath,
		)
		if err := os.RemoveAll(finalPath); err != nil { // Use RemoveAll in case it's a directory
			return fmt.Errorf("failed to remove existing non-executable %s: %w", finalPath, err)
		}
	} else if !os.IsNotExist(err) {
		// An error other than "file does not exist" occurred with os.Stat
		return fmt.Errorf("failed to check status of %s in tools/: %w", finalPath, err)
	}
	// If os.IsNotExist(err) is true, or if we've removed an invalid existing file, proceed with installation.

	fmt.Println("Installing Tailwind CSS to tools/ ...")

	tailwindURL := "https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.7/tailwindcss-linux-x64"
	// downloadPath is also relative to the tools directory
	downloadPath := "tailwindcss-linux-x64"

	// Clean up potentially existing temporary download file
	_ = os.Remove(downloadPath)

	fmt.Printf("Downloading Tailwind CSS from %s to %s (in tools/)...\n", tailwindURL, downloadPath)
	// The run command's workDir is "", so it executes in the current directory (tools/)
	err := run("", "curl", "-sLO", tailwindURL)
	if err != nil {
		return fmt.Errorf("failed to download Tailwind CSS: %w", err)
	}

	// Rename the downloaded file
	fmt.Printf("Renaming %s to %s (in tools/)...\n", downloadPath, finalPath)
	if err := os.Rename(downloadPath, finalPath); err != nil {
		_ = os.Remove(downloadPath) // Attempt to clean up downloaded file if rename fails
		return fmt.Errorf("failed to rename %s to %s: %w", downloadPath, finalPath, err)
	}

	// Make it executable
	fmt.Printf("Making %s executable (in tools/)...\n", finalPath)
	if err := os.Chmod(finalPath, 0o755); err != nil {
		_ = os.Remove(finalPath) // Attempt to clean up renamed file if chmod fails
		return fmt.Errorf("failed to make %s executable: %w", finalPath, err)
	}

	fmt.Printf("Tailwind CSS (%s) installed successfully in tools/.\n", finalPath)
	return nil
}

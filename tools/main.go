//go:build mage

package main

import (
	"errors"
	"fmt"
	"log/slog"
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
		slog.Error(
			"Failed to read go.work file",
			"path",
			workFilePath,
			"cwd",
			wd,
			"original_error",
			err.Error(),
		)
		return nil, errors.New("failed to read go.work file")
	}
	wf, err := modfile.ParseWork(filepath.Base(workFilePath), data, nil)
	if err != nil {
		slog.Error(
			"Failed to parse go.work file",
			"path",
			workFilePath,
			"original_error",
			err.Error(),
		)
		return nil, errors.New("failed to parse go.work file")
	}
	var modules []string
	for _, use := range wf.Use {
		modules = append(modules, use.Path)
	}
	return modules, nil
}

func Lint() error {
	slog.Info("Reading modules from go.work (from tools/magefile.go)")
	modules, err := getWorkspaceModules()
	if err != nil {
		// getWorkspaceModules already logs the error, so we just return a new one.
		return errors.New("could not get workspace modules")
	}

	if len(modules) == 0 {
		slog.Info("No modules found in go.work. Nothing to lint.")
		return nil
	}

	slog.Info("Found modules to lint", "modules", modules)
	slog.Info("Linting Go modules in workspace (from tools/magefile.go)")

	for _, modulePath := range modules {
		slog.Info("Linting module", "module", modulePath)
		modulePath = filepath.Clean(modulePath)

		// This is the actual directory where the module resides, relative to workspace root
		// e.g., "apps/client"
		// We need to construct the path to this directory relative to tools/
		relModuleDir := filepath.Join("..", modulePath) // e.g., ../apps/client

		args := []string{
			"run",
			"--build-tags=mage", // Include files with 'mage' build tag
			"./...",             // Lint all packages within the CWD
		}

		slog.Info("Running golangci-lint", "directory", relModuleDir, "args", args)
		// Call run with the target directory
		err := run(relModuleDir, "golangci-lint", args...)
		if err != nil {
			slog.Error(
				"golangci-lint failed for module",
				"module",
				modulePath,
				"directory",
				relModuleDir,
				"original_error",
				err.Error(),
			)
			return fmt.Errorf(
				"golangci-lint failed for module %s",
				modulePath,
			) // Use fmt.Errorf to include modulePath in the returned error message
		}
		slog.Info("Finished linting module", "module", modulePath)
	}

	slog.Info("All modules linted successfully (by tools/magefile.go)")
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
			slog.Info(
				"Tailwind CSS already installed and executable. Skipping installation.",
				"path",
				finalPath,
			)
			return nil
		}
		// File exists but is not a valid executable (e.g., a directory, or not executable).
		// We'll proceed to remove it and reinstall.
		slog.Warn(
			"Found Tailwind CSS in tools, but it's not a valid executable or has wrong permissions. Proceeding with re-installation.",
			"path",
			finalPath,
		)
		if err := os.RemoveAll(finalPath); err != nil { // Use RemoveAll in case it's a directory
			slog.Error(
				"Failed to remove existing non-executable Tailwind CSS",
				"path",
				finalPath,
				"original_error",
				err.Error(),
			)
			return fmt.Errorf("failed to remove existing non-executable %s", finalPath)
		}
	} else if !os.IsNotExist(err) {
		// An error other than "file does not exist" occurred with os.Stat
		slog.Error("Failed to check status of Tailwind CSS executable", "path", finalPath, "original_error", err.Error())
		return fmt.Errorf("failed to check status of %s", finalPath)
	}
	// If os.IsNotExist(err) is true, or if we've removed an invalid existing file, proceed with installation.

	slog.Info("Installing Tailwind CSS to tools/")

	tailwindURL := "https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.7/tailwindcss-linux-x64"
	// downloadPath is also relative to the tools directory
	downloadPath := "tailwindcss-linux-x64"

	// Clean up potentially existing temporary download file
	_ = os.Remove(downloadPath)

	slog.Info("Downloading Tailwind CSS", "url", tailwindURL, "destination", downloadPath)
	// The run command's workDir is "", so it executes in the current directory (tools/)
	err := run("", "curl", "-sLO", tailwindURL)
	if err != nil {
		slog.Error(
			"Failed to download Tailwind CSS",
			"url",
			tailwindURL,
			"original_error",
			err.Error(),
		)
		return errors.New("failed to download Tailwind CSS")
	}

	// Rename the downloaded file
	slog.Info("Renaming downloaded file", "from", downloadPath, "to", finalPath)
	if err := os.Rename(downloadPath, finalPath); err != nil {
		_ = os.Remove(downloadPath) // Attempt to clean up downloaded file if rename fails
		slog.Error(
			"Failed to rename downloaded Tailwind CSS file",
			"from",
			downloadPath,
			"to",
			finalPath,
			"original_error",
			err.Error(),
		)
		return fmt.Errorf("failed to rename %s to %s", downloadPath, finalPath)
	}

	// Make it executable
	slog.Info("Making file executable", "path", finalPath)
	if err := os.Chmod(finalPath, 0o755); err != nil {
		_ = os.Remove(finalPath) // Attempt to clean up renamed file if chmod fails
		slog.Error(
			"Failed to make Tailwind CSS executable",
			"path",
			finalPath,
			"original_error",
			err.Error(),
		)
		return fmt.Errorf("failed to make %s executable", finalPath)
	}

	slog.Info("Tailwind CSS installed successfully", "path", finalPath)
	return nil
}

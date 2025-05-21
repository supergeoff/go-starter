//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh" // Or your preferred way to run commands
)

// Lint delegates to the magefile in the tools directory
func Lint() error {
	fmt.Println("Delegating lint to tools...")
	// Change directory to tools, then run mage, then change back
	// Or directly invoke mage with -d
	return sh.RunV("mage", "-d", "./tools", "lint")
}

func Serve(dirpath string) error {
	return sh.RunV("mage", "-d", "./tools", "serve", dirpath)
}

func Install() error {
	return sh.RunV("mage", "-d", "./tools", "install")
}

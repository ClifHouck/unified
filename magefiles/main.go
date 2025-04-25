//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func BuildGenerators() error {
	err := sh.Run("go", "build", "-o",
		"./tools/generators/bin/generate_stream_handlers",
		"./tools/generators/generate_stream_handlers.go")
	if err != nil {
		return err
	}
	return nil
}

func GenerateStreamHandlers() error {
	mg.Deps(BuildGenerators)
	err := sh.Run("./tools/generators/bin/generate_stream_handlers")
	if err != nil {
		return err
	}
	return nil
}

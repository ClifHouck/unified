//go:build mage

package main

import (
	"fmt"

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
	fmt.Println("Built generate_stream_handlers")
	return nil
}

func GenerateStreamHandlers() error {
	mg.Deps(BuildGenerators)
	err := sh.Run("./tools/generators/bin/generate_stream_handlers")
	if err != nil {
		return err
	}
	fmt.Println("Ran generate_stream_handlers")
	return nil
}

func Build() error {
	mg.Deps(GenerateStreamHandlers)
	err := sh.Run("go", "build")
	if err != nil {
		return err
	}
	fmt.Println("Built repository")
	return nil
}

func BuildExamples() error {
	mg.Deps(Build)
	err := sh.Run("go", "build", "-o", "examples/bin/doorbell",
		"./examples/main.go")
	if err != nil {
		return err
	}
	fmt.Println("Built examples")
	return nil
}

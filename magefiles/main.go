//go:build mage

package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const GENERATE_STREAM_HANDLERS_BINARY = "./tools/generators/bin/generate_stream_handlers"

// Builds any generators found under './tools/generators'.
func BuildGenerators() error {
	dest := GENERATE_STREAM_HANDLERS_BINARY
	source := "./tools/generators/generate_stream_handlers.go"

	logger := log.WithFields(log.Fields{"destination": dest})

	outOfDate, err := target.Path(dest, source)
	if err != nil {
		return err
	}

	if outOfDate {
		err := sh.Run("go", "build", "-o", dest, source)
		if err != nil {
			return err
		}
		logger.Info("Built generators.")
	} else {
		logger.Info("Generators up to date.")
	}
	return nil
}

// Runs the program which generates the client stream handlers.
func GenerateStreamHandlers() error {
	mg.Deps(BuildGenerators)

	source := GENERATE_STREAM_HANDLERS_BINARY
	destFiles := []string{
		"./client/protect_device_update_stream_handler.go",
		"./client/protect_event_stream_handler.go",
	}

	logger := log.WithFields(log.Fields{
		"source":      source,
		"destination": "./client/*_stream_handler.go"})

	var outOfDate bool
	for _, dest := range destFiles {
		changed, err := target.Path(dest, source)
		if err != nil {
			return err
		}
		outOfDate = outOfDate || changed
	}

	if outOfDate {
		err := sh.Run(source)
		if err != nil {
			return err
		}

		err = sh.Run("go", "build", "./client/")
		if err != nil {
			return err
		}

		logger.Info("Stream handlers regenerated.")
	} else {
		logger.Info("Stream handlers up to date.")
	}
	return nil
}

// Builds core *.go files of this respository.
func Build() error {
	mg.Deps(GenerateStreamHandlers)

	err := os.Mkdir("build", 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	dest := "build/build.marker"

	logger := log.WithFields(log.Fields{
		"destination": dest,
	})

	outOfDate, err := target.Glob(dest,
		"./client/*.go",
		"./types/*.go",
	)
	if err != nil {
		return err
	}

	if outOfDate {
		err = sh.Run("go", "build", "./...")
		if err != nil {
			return err
		}
		logger.Info("Built core module go files.")

		err = sh.Run("touch", dest)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Core go build up to date.")
	}
	return nil
}

// Builds the unified CLI command.
func BuildCmd() error {
	mg.Deps(Build)

	dest := "build/unified"

	logger := log.WithFields(log.Fields{
		"destination": dest,
	})

	// TODO: This surprised me... because even though build was triggered, this
	// did not rebuild. Which makes sense in retro
	outOfDate, err := target.Glob(dest,
		"./cmd/*.go",
		"./client/*.go",
		"./types/*.go",
	)
	if err != nil {
		return err
	}

	if outOfDate {
		err := sh.Run("go", "build", "-o", "build/unified",

			"./main.go")
		if err != nil {
			return err
		}
		logger.Info("Built 'unified' cmd.")
	} else {
		logger.Info("'unified' up to date.")
	}

	return nil
}

// Builds example programs found in './examples'.
func BuildExamples() error {
	mg.Deps(Build)
	err := sh.Run("go", "build", "-o", "examples/bin/doorbell",
		"./examples/main.go")
	if err != nil {
		return err
	}
	log.Info("Built examples.")
	return nil
}

// Runs tests.
func Test() error {
	mg.Deps(GenerateStreamHandlers)
	err := sh.RunV("go", "test", "./...")
	if err != nil {
		return err
	}
	log.Info("Ran tests.")
	return nil
}

// Cleans up any built files.
func Clean() error {
	files := []string{
		GENERATE_STREAM_HANDLERS_BINARY,
		// TODO: Should these be checked-in or should they always be re-generated?
		"./client/protect_device_update_stream_handler.go",
		"./client/protect_event_stream_handler.go",
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return err
		}
		log.Infof("Removed '%s'", file)
	}

	paths := []string{
		"build/",
		"examples/bin",
		"tools/generators/bin",
	}
	for _, path := range paths {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
		log.Infof("Removed path '%s'", path)
	}
	return nil
}

// Run the linter
func Lint() error {
	err := sh.RunV("golangci-lint", "run")
	if err != nil {
		return err
	}
	log.Info("Ran linter.")
	return nil
}

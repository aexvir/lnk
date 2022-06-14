//go:build mage

package main

// notest // task orchestrator, not part of the service code

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/princjef/mageutil/bintool"
)

// Build the lnk server binary
func Build(ctx context.Context) error {
	return run(ctx, "go", "build", "-o", "bin/lnkctl", ".")
}

// Run the lnk server without building it
func Run(ctx context.Context) error {
	return run(ctx, "go", "run", ".")
}

// Generate proto definitions and code
func Generate(ctx context.Context) error {
	return run(ctx, "go", "generate", "./...")
}

// Format codebase using gofmt, goimports and prettier
func Format(ctx context.Context) error {
	errfmt := run(ctx, "gofmt", "-w", "-s", ".")
	if errfmt != nil {
		color.Red("failed to format code: %s", errfmt.Error())
	}

	errimp := run(ctx, "goimports", "-w", "-local", "github.com/aexvir/lnk", ".")
	if errimp != nil {
		color.Red("failed to format code: %s", errimp.Error())
	}

	color.Green("looking sweet!")
	return nil
}

// lint the code using go mod tidy and golangci-lint
func Lint(ctx context.Context) error {
	// run go mod tidy and check differences
	gomod, _ := os.ReadFile("go.mod")
	gosum, _ := os.ReadFile("go.sum")

	err := run(ctx, "go", "mod", "tidy", "-v")
	if err != nil {
		return err
	}

	newmod, _ := os.ReadFile("go.mod")
	newsum, _ := os.ReadFile("go.sum")

	if !bytes.Equal(gomod, newmod) || !bytes.Equal(gosum, newsum) {
		return errors.New("differences found; fixed go module")
	}

	color.Green("no differences found")

	// run golangci-lint
	gci, _ := bintool.New(
		"golangci-lint{{.BinExt}}",
		"1.46.2",
		"https://github.com/golangci/golangci-lint/releases/download/v{{.Version}}/golangci-lint-{{.Version}}-{{.GOOS}}-{{.GOARCH}}{{.ArchiveExt}}",
	)
	gci.Ensure()

	return gci.Command("run --max-same-issues 0 --max-issues-per-linter 0").Run()
}

// Test the whole codebased
func Test(ctx context.Context) error {
	return run(ctx, "go", "test", "-race", "-cover", "./...")
}

// cmd builds a command where the stdout and stderr...
func cmd(ctx context.Context, program string, args ...string) *exec.Cmd {
	command := exec.CommandContext(ctx, program, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command
}

// run a command...
func run(ctx context.Context, program string, args ...string) error {
	command := cmd(ctx, program, args...)
	fmt.Printf(
		"%s %s\n",
		color.MagentaString(">"),
		color.New(color.Bold).Sprint(program, " ", strings.Join(args, " ")),
	)

	return command.Run()
}

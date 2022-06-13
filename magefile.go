//go:build mage

package main

// notest // task orchestrator, not part of the service code

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// Build the lnk server binary
func Build(ctx context.Context) error {
	return run(ctx, "go", "build", "-o", "bin/lnkctl", ".")
}

// Run the lnk server
func Run(ctx context.Context) error {
	return run(ctx, "go", "run", ".")
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

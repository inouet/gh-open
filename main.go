package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/skratchdot/open-golang/open"
	"os"
	"strings"
)

type options struct {
	Line     string `short:"l" long:"line" description:"Line number (10 or 10-20)"`
	Branch   string `short:"b" long:"branch" description:"Branch name"`
	PrintURL bool   `short:"p" long:"print" description:"Print url"`
}

var (
	opts        options
	statusOK    = 0
	statusError = 1
)

func printError(err error) {
	fmt.Printf("Error: %+v\n", err)
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	parser := flags.NewParser(&opts, flags.Default)
	args, err := parser.Parse()

	if err != nil {
		printError(err)
		return statusError
	}

	if len(args) == 0 {
		parser.WriteHelp(os.Stdout)
		return statusError
	}

	objectPath := args[0]

	line := strings.TrimSpace(opts.Line)
	if !validateLine(line) {
		printError(fmt.Errorf("invalid line format"))
		return statusError
	}

	gr, err := newGitRemote(objectPath)

	if err != nil {
		printError(err)
		return statusError
	}

	remoteURL, err := gr.remoteURL(opts.Branch, line)
	if err != nil {
		printError(err)
		return statusError
	}

	if opts.PrintURL {
		fmt.Println(remoteURL)
		return statusOK
	}

	err = open.Run(remoteURL)
	if err != nil {
		printError(err)
		return statusError
	}
	return statusOK
}

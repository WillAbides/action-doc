package main

import (
	"fmt"
	"io"
	"os"

	actiondoc "github.com/willabides/action-doc"
)

var osExit = os.Exit

func main() {
	run(os.Stdin, os.Stdout, os.Stderr, osExit)
	err := os.Stdout.Close()
	if err != nil {
		panic(err)
	}
}

func run(stdin io.Reader, stdout, stderr io.Writer, osExit func(int)) {
	markdown, err := actiondoc.ActionMarkdown(stdin)
	if err != nil {
		fmt.Fprintf(stderr, "error reading action definition: %v", err)
		osExit(1)
		return
	}
	_, err = stdout.Write(markdown)
	if err != nil {
		fmt.Fprintf(stderr, "error writing: %v", err)
		osExit(1)
	}
}

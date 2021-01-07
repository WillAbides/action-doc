package main

import (
	"fmt"
	"os"

	actiondoc "github.com/willabides/action-doc"
)

var osExit = os.Exit

func exitOnErr(err error, stmt string, args ...interface{}) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, stmt, args...)
	osExit(1)
}

func main() {
	markdown, err := actiondoc.ActionMarkdown(os.Stdin)
	exitOnErr(err, "error reading action definition: %v", err)
	_, err = os.Stdout.Write(markdown)
	exitOnErr(err, "error writing: %v", err)
	err = os.Stdout.Close()
	exitOnErr(err, "error flushing stdout: %v", err)
}

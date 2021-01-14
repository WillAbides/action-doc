package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"
	actiondoc "github.com/willabides/action-doc"
)

type cliOptions struct {
	Version kong.VersionFlag `help:"Show version and exit"`

	SkipActionAuthor      bool   `help:"Skip outputting the action author"`
	SkipActionName        bool   `help:"Skip outputting the action name"`
	SkipActionDescription bool   `help:"Skip outputting the action description"`
	PostDescriptionText   string `help:"Some text to insert after the description"`
	HeaderPrefix          string `help:"Some extra #s for the markdown headers"`
	ActionConfig          string `kong:"arg,help='action.yml to parse'"`
}

var cli cliOptions

func run(stdout io.Writer) error {
	opts := []actiondoc.MarkdownOption{
		actiondoc.PostDescriptionText(cli.PostDescriptionText),
		actiondoc.HeaderPrefix(cli.HeaderPrefix),
		actiondoc.SkipActionName(cli.SkipActionName),
		actiondoc.SkipActionDescription(cli.SkipActionDescription),
		actiondoc.SkipActionAuthor(cli.SkipActionAuthor),
	}
	input, err := os.Open(cli.ActionConfig)
	if err != nil {
		return err
	}
	markdown, err := actiondoc.ActionMarkdown(input, opts...)
	if err != nil {
		return fmt.Errorf("error parsing action definition: %v", err)
	}
	_, err = stdout.Write(markdown)
	if err != nil {
		return fmt.Errorf("error writing: %v", err)
	}
	return err
}

var version = "unknown"

func main() {
	kong.Parse(&cli, kong.Vars{
		"version": version,
	}).FatalIfErrorf(run(os.Stdout))
}

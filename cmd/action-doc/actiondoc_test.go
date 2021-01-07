package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errWriter struct{}

func (b *errWriter) Write(p []byte) (int, error) {
	return 0, assert.AnError
}

func testdataActions(filename string) string {
	return filepath.Join(filepath.FromSlash("../../testdata/actions"), filename)
}

func Test_run(t *testing.T) {
	for _, td := range []struct {
		name     string
		wantFile string
		cli      cliOptions
	}{
		{
			name:     "as it comes",
			wantFile: "ex1.md",
		},
		{
			name:     "--skip-action-name",
			wantFile: "ex1-skip_name.md",
			cli: cliOptions{
				SkipActionName: true,
			},
		},
		{
			name:     "--skip-action-author",
			wantFile: "ex1-skip_author.md",
			cli: cliOptions{
				SkipActionAuthor: true,
			},
		},
		{
			name:     "--skip-action-description",
			wantFile: "ex1-skip_action_description.md",
			cli: cliOptions{
				SkipActionDescription: true,
			},
		},
		{
			name:     "gotta skip 'em all",
			wantFile: "ex1-skip_all.md",
			cli: cliOptions{
				SkipActionDescription: true,
				SkipActionAuthor:      true,
				SkipActionName:        true,
			},
		},
	} {
		t.Run(td.name, func(t *testing.T) {
			want, err := ioutil.ReadFile(testdataActions(td.wantFile))
			require.NoError(t, err)
			cli = td.cli
			if cli.ActionConfig == "" {
				cli.ActionConfig = testdataActions("ex1.yml")
			}
			var stdout bytes.Buffer
			err = run(&stdout)
			require.NoError(t, err)
			require.Equal(t, string(want), stdout.String())
		})
	}

	t.Run("invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		actionConfig := filepath.Join(tmpDir, "action.yml")
		err := ioutil.WriteFile(actionConfig, []byte("boogers"), 0o600)
		require.NoError(t, err)
		cli = cliOptions{
			ActionConfig: actionConfig,
		}
		var stdout bytes.Buffer
		err = run(&stdout)
		require.Error(t, err)
		require.True(t, strings.HasPrefix(err.Error(), "error parsing action definition:"))
		require.Empty(t, stdout.String())
	})

	t.Run("missing file", func(t *testing.T) {
		cli = cliOptions{
			ActionConfig: "fake_file.yml",
		}
		var stdout bytes.Buffer
		err := run(&stdout)
		require.EqualError(t, err, "open fake_file.yml: no such file or directory")
		require.Empty(t, stdout.String())
	})

	t.Run("bad stdout", func(t *testing.T) {
		cli = cliOptions{
			ActionConfig: testdataActions("ex1.yml"),
		}
		err := run(new(errWriter))
		require.EqualError(t, err, "error writing: assert.AnError general error for testing")
	})
}

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setOSStdin(t *testing.T, file *os.File) {
	t.Helper()
	stdin := os.Stdin
	os.Stdin = file
	t.Cleanup(func() {
		os.Stdin = stdin
	})
}

func captureOsOut(t *testing.T, target string, fn func()) []byte {
	t.Helper()
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, target)
	writer, err := os.Create(filename)
	require.NoError(t, err)
	switch target {
	case "stdout":
		orig := os.Stdout
		os.Stdout = writer
		fn()
		os.Stdout = orig
	case "stderr":
		orig := os.Stderr
		os.Stderr = writer
		fn()
		os.Stderr = orig
	default:
		panic("must be stdout or stderr")
	}
	got, err := ioutil.ReadFile(filename)
	require.NoError(t, err)
	return got
}

func captureAllOSOut(t *testing.T, fn func()) (stdout, stderr []byte) {
	t.Helper()
	stdout = captureOsOut(t, "stdout", func() {
		stderr = captureOsOut(t, "stderr", fn)
	})
	return stdout, stderr
}

func requireExit(t *testing.T, exits bool, code int) {
	exited := false
	exitCode := 0
	orig := osExit
	osExit = func(n int) {
		if exited {
			t.Error("only one exit allowed")
		}
		exited = true
		exitCode = n
	}
	t.Cleanup(func() {
		osExit = orig
		assert.Equal(t, exits, exited)
		assert.Equal(t, code, exitCode)
	})
}

func Test_main(t *testing.T) {
	t.Run("invalid yaml", func(t *testing.T) {
		requireExit(t, true, 1)
		tmpDir := t.TempDir()
		err := ioutil.WriteFile(filepath.Join(tmpDir, "stdin"), []byte("boogers"), 0o600)
		require.NoError(t, err)
		stdin, err := os.Open(filepath.Join(tmpDir, "stdin"))
		require.NoError(t, err)
		setOSStdin(t, stdin)
		stdout, stderr := captureAllOSOut(t, main)
		require.Empty(t, stdout)
		require.True(t, strings.HasPrefix(string(stderr), "error reading action definition:"))
	})

	t.Run("happy path", func(t *testing.T) {
		requireExit(t, false, 0)
		stdin, err := os.Open(filepath.FromSlash("../../testdata/actions/ex1.yml"))
		require.NoError(t, err)
		setOSStdin(t, stdin)
		want, err := ioutil.ReadFile(filepath.FromSlash("../../testdata/actions/ex1.md"))
		require.NoError(t, err)
		var stderr []byte
		stdout := captureOsOut(t, "stdout", func() {
			stderr = captureOsOut(t, "stderr", func() {
				main()
			})
		})
		require.Equal(t, string(want), string(stdout))
		require.Empty(t, stderr)
	})
}

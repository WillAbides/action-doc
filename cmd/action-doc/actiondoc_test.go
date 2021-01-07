package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

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

func Test_main(t *testing.T) {
	exited := false
	osExit = func(_ int) {
		exited = true
	}
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
	require.False(t, exited)
}

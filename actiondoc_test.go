package actiondoc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActionMarkdown(t *testing.T) {
	actionFile, err := os.Open(filepath.FromSlash("testdata/actions/ex1.yml"))
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, actionFile.Close())
	})
	want, err := ioutil.ReadFile(filepath.FromSlash("testdata/actions/ex1.md"))
	require.NoError(t, err)
	got, err := ActionMarkdown(actionFile)
	require.NoError(t, err)
	require.Equal(t, string(want), string(got))
}

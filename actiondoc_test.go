package actiondoc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActionMarkdown(t *testing.T) {
	for _, td := range []struct {
		name     string
		opts     []MarkdownOption
		wantFile string
	}{
		{
			name:     "as it comes",
			wantFile: "testdata/actions/ex1.md",
		},
		{
			name:     "SkipActionName",
			opts:     []MarkdownOption{SkipActionName(true)},
			wantFile: "testdata/actions/ex1-skip_name.md",
		},
		{
			name:     "SkipActionDescription",
			opts:     []MarkdownOption{SkipActionDescription(true)},
			wantFile: "testdata/actions/ex1-skip_action_description.md",
		},
		{
			name:     "SkipActionAuthor",
			opts:     []MarkdownOption{SkipActionAuthor(true)},
			wantFile: "testdata/actions/ex1-skip_author.md",
		},
		{
			name:     "HeaderPrefix",
			opts:     []MarkdownOption{HeaderPrefix("##")},
			wantFile: "testdata/actions/ex1-header_prefix.md",
		},
		{
			name: "skip all",
			opts: []MarkdownOption{
				SkipActionAuthor(true),
				SkipActionDescription(true),
				SkipActionName(true),
			},
			wantFile: "testdata/actions/ex1-skip_all.md",
		},
	} {
		t.Run(td.name, func(t *testing.T) {
			actionFile, err := os.Open(filepath.FromSlash("testdata/actions/ex1.yml"))
			require.NoError(t, err)
			t.Cleanup(func() {
				require.NoError(t, actionFile.Close())
			})
			var got []byte
			got, err = ActionMarkdown(actionFile, td.opts...)
			require.NoError(t, err)
			fmt.Println(string(got))
			var want []byte
			want, err = ioutil.ReadFile(filepath.FromSlash(td.wantFile))
			require.NoError(t, err)
			require.Equal(t, string(want), string(got))
		})
	}
}

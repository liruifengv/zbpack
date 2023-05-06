package source_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeabur/zbpack/internal/source"
)

func getGithubToken(t *testing.T) string {
	token, ok := os.LookupEnv("GITHUB_TOKEN")

	if !ok {
		t.Skip("no token (GITHUB_TOKEN) provided: skipping GitHub tests")
	}

	return token
}

func TestGitHubFsOpen_File(t *testing.T) {
	token := getGithubToken(t)

	fs := source.NewGitHubFs("zeabur", "zeabur", token)
	f, err := fs.Open("readme.md")
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestGitHubFsOpen_Dir(t *testing.T) {
	token := getGithubToken(t)

	fs := source.NewGitHubFs("zeabur", "zeabur", token)
	f, err := fs.Open("")
	assert.NoError(t, err)
	assert.NotNil(t, f)
}

func TestGitHubFsOpenFile_WithWriteFlag(t *testing.T) {
	token := getGithubToken(t)

	fs := source.NewGitHubFs("zeabur", "zeabur", token)
	_, err := fs.OpenFile("readme.md", os.O_RDWR, 0)
	assert.ErrorAs(t, err, source.ErrReadonly)
}
package gitutils

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestGetUpstream(t *testing.T) {
	config := `
[branch "master"]
	remote = origin
	merge = refs/heads/master
[branch "feature/projects"]
	remote = spooky
	merge = refs/heads/feature/oldprojects
`

	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("ref: refs/heads/master\n"),
		},
		".git/config": &fstest.MapFile{
			Data: []byte(config),
		},
	}

	git := &GitUtils{
		pathToGit: "git",
		fsys:      files,
		RepoRoot:  "/Users/oriana/dev/kitsch-prompt",
	}

	assert.Equal(t,
		"origin/master",
		git.GetUpstream("master"),
	)

	assert.Equal(t,
		"spooky/feature/oldprojects",
		git.GetUpstream("feature/projects"),
	)

	assert.Equal(t,
		"",
		git.GetUpstream("banana"),
	)
}

func TestGetUpstreamNoConfig(t *testing.T) {
	files := fstest.MapFS{
		".git/HEAD": &fstest.MapFile{
			Data: []byte("ref: refs/heads/master\n"),
		},
	}

	git := &GitUtils{
		pathToGit: "git",
		fsys:      files,
		RepoRoot:  "/Users/oriana/dev/kitsch-prompt",
	}

	assert.Equal(t,
		"",
		git.GetUpstream("feature/projects"),
	)
}

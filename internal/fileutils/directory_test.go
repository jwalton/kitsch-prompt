package fileutils

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
)

func TestHasExtension(t *testing.T) {
	fsys := fstest.MapFS{
		"foo.go": &fstest.MapFile{
			Data: []byte(""),
		},
		"foo.test.js": &fstest.MapFile{
			Data: []byte(""),
		},
		".gitignore": &fstest.MapFile{
			Data: []byte(""),
		},
		"banana": &fstest.MapFile{
			Data: []byte(""),
		},
	}

	dir := NewDirectoryTestFS("/foo/bar", fsys)

	assert.Equal(t, true, dir.HasExtension("go"))
	assert.Equal(t, true, dir.HasExtension(".go"))
	assert.Equal(t, true, dir.HasExtension("js"))
	assert.Equal(t, true, dir.HasExtension(".js"))
	assert.Equal(t, true, dir.HasExtension("test.js"))
	assert.Equal(t, true, dir.HasExtension(".test.js"))

	assert.Equal(t, false, dir.HasExtension("banana"))
	assert.Equal(t, false, dir.HasExtension("gitignore"))
	assert.Equal(t, false, dir.HasExtension(".gitignore"))
	assert.Equal(t, false, dir.HasExtension("txt"))
	assert.Equal(t, false, dir.HasExtension(".txt"))
}

func TestHasFile(t *testing.T) {
	fsys := fstest.MapFS{
		"src/index.js": &fstest.MapFile{
			Data: []byte(""),
		},
	}

	dir := NewDirectoryTestFS("/foo/bar", fsys)

	assert.Equal(t, true, dir.HasFile("src"))
	assert.Equal(t, true, dir.HasFile("src/index.js"))
	assert.Equal(t, false, dir.HasDirectory("test"))
}

func TestHasDirectory(t *testing.T) {
	fsys := fstest.MapFS{
		"src/index.js": &fstest.MapFile{
			Data: []byte(""),
		},
	}

	dir := NewDirectoryTestFS("/foo/bar", fsys)

	assert.Equal(t, true, dir.HasDirectory("src"))
	assert.Equal(t, false, dir.HasDirectory("test"))
}

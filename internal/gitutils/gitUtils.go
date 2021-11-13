package gitutils

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
)

// ErrNoGit is emitted when we need to run the git executable, but git is not installed.
var ErrNoGit = errors.New("Git is not installed")

// GitUtils is an object that allows you to retrieve information about
// a git repository.
type GitUtils struct {
	// pathToGit is the path to the git executable.
	pathToGit string
	// fsys is an fs.FS instance bound to the root of the git repository.
	fsys fs.FS
	// RepoRoot is the root folder of the git repository.
	RepoRoot string
}

// New returns a new instance of `GitUtils` for the specified repository.
func New(pathToGit string, folder string) *GitUtils {
	// Resolve the path to the git executable
	pathToGit, err := fileutils.LookPathSafe(pathToGit)
	if err != nil {
		pathToGit = ""
	}

	// Figure out whether or not we're inside a git repo.
	gitRoot := FindGitRoot(folder)

	var fsys fs.FS = nil
	if gitRoot != "" {
		fsys = os.DirFS(gitRoot)
	}

	if gitRoot == "" {
		return nil
	}

	return &GitUtils{
		pathToGit: pathToGit,
		fsys:      fsys,
		RepoRoot:  gitRoot,
	}
}

// FindGitRoot returns the root of the current git repo.
func FindGitRoot(cwd string) string {
	gitFolder := fileutils.FindFileInAncestors(cwd, ".git")
	if gitFolder != "" {
		return filepath.Dir(gitFolder)
	}
	return ""
}

// git will run a git command in the root folder of the git repository.
// Returns empty string if there was an error running the command.
func (utils *GitUtils) git(args ...string) (string, error) {
	if utils.pathToGit == "" {
		return "", ErrNoGit
	}

	cmd := exec.Command(utils.pathToGit, args...)
	cmd.Dir = utils.RepoRoot

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func countLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0

	for {
		n, err := r.Read(buf)
		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				count++
			}
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
	}

	return count, nil
}

// GetStashCount returns the number of stashes, or 0 if there are none or
// the path is not a git repo.
//
// `path` should be the git root folder.
func (utils *GitUtils) GetStashCount() (int, error) {
	// TODO: Read .git/logs/refs/stash, and count the number of `\n`s.`
	file, err := utils.fsys.Open(".git/logs/refs/stash")
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}

	defer file.Close()
	return countLines(file)
}

// // GetCurrentRepo returns a git repo for the current folder, or nil if we are not
// // inside a git repo.
// func OpenRepo(path string) *git.Repository {
// 	gitFolder := fileutils.FindFileInAncestors(path, ".git")

// 	repo, err := git.PlainOpen(gitFolder)

// 	if err != nil {
// 		return nil
// 	}

// 	return repo
// }

// // GetShortName returns the short name for the given reference.  This will
// // be the branch name, the tag name, or the hash.
// func GetShortName(repo *git.Repository, ref *plumbing.Reference) string {
// 	var shortName string

// 	// If this is a branch, return the branch name
// 	refName := ref.Name()
// 	if refName.IsBranch() {
// 		shortName = refName.Short()
// 	}

// 	if shortName == "" {
// 		// Search for a tag with this ref.
// 		shortName = getTagName(repo, ref)
// 	}

// 	if shortName == "" {
// 		// If all else fails, use the hash.
// 		shortName = "(" + ref.Hash().String()[0:7] + "...)"
// 	}

// 	return shortName
// }

// func getTagName(repo *git.Repository, ref *plumbing.Reference) string {
// 	var result string

// 	if ref.Name().IsTag() {
// 		result = ref.Name().Short()
// 	} else {
// 		annotatedTag, _ := repo.TagObject(ref.Hash())
// 		if annotatedTag != nil {
// 			result = annotatedTag.Name
// 		} else {
// 			// Need to search for the tag.
// 			tags, err := repo.Tags()
// 			CheckIfError(err)
// 			err = tags.ForEach(func(t *plumbing.Reference) error {
// 				if t.Hash() == ref.Hash() {
// 					result = t.Name().Short()
// 					return storer.ErrStop
// 				}
// 				return nil
// 			})
// 		}
// 	}

// 	if result != "" {
// 		result = "(" + result + ")"
// 	}

// 	return result
// }

package gitutils

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jwalton/kitsch-prompt/internal/fileutils"
)

// GitUtils is an object that allows you to retrieve information about
// a git repository.
type GitUtils struct {
	pathToGit string
	// RepoRoot is the root folder of the git repository.
	RepoRoot string
}

// New returns a new instance of `GitUtils` for the specified repository.
func New(pathToGit string, path string) *GitUtils {
	gitRoot := findGitRoot(pathToGit, path)

	if gitRoot == "" {
		return nil
	}

	return &GitUtils{
		pathToGit,
		gitRoot,
	}
}

// FindGitRoot returns the root of the current git repo.
func FindGitRoot(cwd string) string {
	return findGitRoot("git", cwd)
}

func findGitRoot(pathToGit string, cwd string) string {
	gitFolder := fileutils.FindFileInAncestors(cwd, ".git")
	if gitFolder != "" {
		return filepath.Dir(gitFolder)
	}
	return ""
}

// git will run a git command in the root folder of the git repository.
// Returns empty string if there was an error running the command.
func (utils *GitUtils) git(args ...string) (string, error) {
	cmd := exec.Command(utils.pathToGit, args...)
	cmd.Dir = utils.RepoRoot
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

// GetStashCount returns the number of stashes, or 0 if there are none or
// the path is not a git repo.
//
// `path` should be the git root folder.
func (utils *GitUtils) GetStashCount() int {
	countStr, err := utils.git("rev-list", "--walk-reflogs", "--count", "refs/stash")
	if err != nil {
		return 0
	}

	result, err := strconv.Atoi(strings.TrimSpace(countStr))
	if err != nil {
		// ???
		return 0
	}

	return result
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
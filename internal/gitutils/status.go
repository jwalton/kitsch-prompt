package gitutils

import "os/exec"

// Stats returns status counters for the given git repo.
func (utils *GitUtils) Stats() (GitStats, error) {
	// This uses `exec.Command` instead of go-git's worktree.Status(),
	// because worktree.Status() is crazy slow: https://github.com/go-git/go-git/issues/181
	cmd := exec.Command(utils.pathToGit, "status", "-z")
	cmd.Dir = utils.RepoRoot
	stats := GitStats{}
	cmd.Stdout = &statusWriter{stats: &stats}
	err := cmd.Run()
	return stats, err
}

// GitStats represents counts about files which are in the index, in the work tree,
// and files which are unmerged.
type GitStats struct {
	// Index contains counts of files in the index.
	Index GitFileStats
	// Files contains counts of unstaged changes in the work tree.
	Files GitFileStats
	// Unmerged is a count of unmerged files.
	Unmerged int
}

// GitFileStats contains counts of files in the index or in the work tree.
type GitFileStats struct {
	// Added is the number of files that have been added.
	Added int
	// Modified is the number of files that have been modified.
	Modified int
	// Deleted is the number of files that have been deleted.
	Deleted int
}

type statusWriter struct {
	linePos int
	stats   *GitStats
}

func countStats(stats *GitFileStats, x byte) {
	switch x {
	case 'M':
		stats.Modified++
	case 'A':
		stats.Added++
	case 'D':
		stats.Deleted++
	case 'R':
		stats.Modified++
	case 'C':
		stats.Modified++
	}
}

// Write parses the output of `git status -z` and counts files in a GitStats.
func (status *statusWriter) Write(p []byte) (n int, err error) {
	var i int
	var x byte

	for i = 0; i < len(p); i++ {
		if status.linePos == 0 {
			x = p[i]
			status.linePos++
		} else if status.linePos == 1 {
			y := p[i]

			if (x == 'D' && y == 'D') || (x == 'A' && y == 'A') || x == 'U' || y == 'U' {
				status.stats.Unmerged++
			} else if x == '?' {
				status.stats.Files.Added++
			} else {
				countStats(&status.stats.Index, x)
				countStats(&status.stats.Files, y)
			}

			status.linePos++
		} else if p[i] == 0 {
			status.linePos = 0
		} else {
			status.linePos++
		}

	}
	return len(p), nil
}
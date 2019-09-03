package git

import "strings"

type ErrOpeningRepository struct {
	Reasons []string
}

func (e ErrOpeningRepository) Error() string {
	return "error opening repository: " + strings.Join(e.Reasons, ", ")
}

type ErrGettingLatestTag struct {
	Reasons []string
}

func (e ErrGettingLatestTag) Error() string {
	return "error getting latest tag: " + strings.Join(e.Reasons, ", ")
}

type ErrGettingLatestCommit struct {
	Reasons []string
}

func (e ErrGettingLatestCommit) Error() string {
	return "error getting latest commit: " + strings.Join(e.Reasons, ", ")
}

type ErrNoTags struct{}

func (e ErrNoTags) Error() string {
	return "no tags"
}

type ErrGeneratingWorkTree struct{}

func (e ErrGeneratingWorkTree) Error() string {
	return "error generating git repo worktree"
}

type ErrGitAdd struct{}

func (e ErrGitAdd) Error() string {
	return "git add error"
}

type ErrGitCommit struct{}

func (e ErrGitCommit) Error() string {
	return "git commit error"
}

type ErrGitPush struct{}

func (e ErrGitPush) Error() string {
	return "git push error"
}
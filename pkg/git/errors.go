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

type ErrGeneratingWorkTree struct{
	Reasons []string
}

func (e ErrGeneratingWorkTree) Error() string {
	return "error generating git repo worktree: " + strings.Join(e.Reasons, ", ")
}

type ErrGitAdd struct{
	Reasons []string
}

func (e ErrGitAdd) Error() string {
	return "git add error: " + strings.Join(e.Reasons, ", ")
}

type ErrGitCommit struct{
	Reasons []string
}

func (e ErrGitCommit) Error() string {
	return "git commit error: " + strings.Join(e.Reasons, ", ")
}

type ErrGitPush struct{
	Reasons []string
}

func (e ErrGitPush) Error() string {
	return "git push error: " + strings.Join(e.Reasons, ", ")
}

type ErrGeneratingRelFilePath struct{
	Reasons []string
}

func (e ErrGeneratingRelFilePath) Error() string {
	return "error generating relative file path: " + strings.Join(e.Reasons, ", ")
}
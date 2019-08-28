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

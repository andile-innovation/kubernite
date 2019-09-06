package git

import (
	"fmt"
	goGit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	//gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Event string

func (e Event) String() string {
	return string(e)
}

const TagEvent Event = "tag"

type Repository struct {
	goGit.Repository
}

func NewRepositoryFromFilePath(pathToRepository string) (*Repository, error) {
	// validate given repository path
	if pathToRepository == "" {
		return nil, ErrOpeningRepository{Reasons: []string{
			"given path to repository is blank",
		}}
	}
	pathToRepository, err := filepath.Abs(pathToRepository)
	if err != nil {
		return nil, ErrOpeningRepository{Reasons: []string{
			"getting absolute path",
			err.Error(),
		}}
	}
	repositoryFileInfo, err := os.Stat(pathToRepository)
	if err != nil {
		return nil, ErrOpeningRepository{Reasons: []string{
			"getting repository file info",
			err.Error(),
		}}
	}
	if !repositoryFileInfo.IsDir() {
		return nil, ErrOpeningRepository{Reasons: []string{
			fmt.Sprintf("'%s' is not a directory", pathToRepository),
		}}
	}

	// try and open repository at given path
	repository, err := goGit.PlainOpen(pathToRepository)
	if err != nil {
		return nil, ErrOpeningRepository{Reasons: []string{
			err.Error(),
		}}
	}

	return &Repository{
		Repository: *repository,
	}, nil
}

func (r *Repository) GetLatestTagName() (string, error) {
	latestTag, err := r.GetLatestTag()
	if err != nil {
		return "", err
	}
	return latestTag.Name().Short(), nil
}

func (r *Repository) GetLatestTag() (*plumbing.Reference, error) {
	tagReferenceIterator, err := r.Tags()
	if err != nil {
		return nil, ErrGettingLatestTag{Reasons: []string{
			err.Error(),
		}}
	}
	if tagReferenceIterator == nil {
		return nil, ErrNoTags{}
	}

	latestTag, err := tagReferenceIterator.Next()
	switch err {
	case io.EOF:
		// no tags
		return nil, ErrNoTags{}
	case nil:
		return latestTag, nil
	default:
		return nil, ErrGettingLatestTag{Reasons: []string{
			"getting next tag",
			err.Error(),
		}}
	}
}

func (r *Repository) GetLatestCommitHash() (string, error) {
	latestCommit, err := r.GetLatestCommit()
	if err != nil {
		return "", err
	}
	return latestCommit.Hash.String(), nil
}

func (r *Repository) GetLatestCommit() (*object.Commit, error) {
	commitReferenceIterator, err := r.CommitObjects()
	if err != nil {
		return nil, ErrGettingLatestCommit{Reasons: []string{
			"getting commit iterator",
			err.Error(),
		}}
	}
	commit, err := commitReferenceIterator.Next()
	if err != nil {
		return nil, ErrGettingLatestCommit{Reasons: []string{
			"getting next commit",
			err.Error(),
		}}
	}
	return commit, nil
}

func (r *Repository) CommitDeployment(DeploymentFileRepositoryPath, DeploymentFilePath, GitUsername, GitPassword, GitKey string) error {
	// get worktree
	w, err := r.Worktree()
	if err != nil {
		return ErrGeneratingWorkTree{
			Reasons: []string{
				"getting git deployment repo worktree",
				err.Error(),
			},
		}
	}

	fileRelToRepo, err := filepath.Rel(DeploymentFileRepositoryPath, DeploymentFilePath)
	if err != nil {
		return ErrGeneratingRelFilePath{
			Reasons: []string{
				"generating deployment file path relative to deployment file repository path",
				err.Error(),
			},
		}
	}

	// git add deploymentFile
	if _, err := w.Add(fileRelToRepo); err != nil {
		return ErrGitAdd{
			Reasons: []string{
				"git add deployment",
				err.Error(),
			},
		}
	}

	// git commit kubernite deployment
	if _, err := w.Commit("Kubernite deployment", &goGit.CommitOptions{
		Author: &object.Signature{
			Name:  "Kubernite",
			Email: "-",
			When:  time.Now(),
		},
	}); err != nil {
		return ErrGitCommit{Reasons: []string{
			"git commit deployment",
			err.Error(),
		}}
	}

	//signer, err := ssh.ParsePrivateKey([]byte(GitKey))
	//auth := &gitssh.PublicKeys{User: "git", Signer: signer}
	//
	//// git push kubernite deployment
	//if err := r.Push(&goGit.PushOptions{
	//	RemoteName: "origin",
	//	RefSpecs:   nil,
	//	Auth:       auth,
	//	Progress:   nil,
	//	Prune:      false,
	//}); err != nil {
	//	return ErrGitPush{
	//		Reasons: []string{
	//			"git push deployment",
	//			err.Error(),
	//		},
	//	}
	//}

	return nil
}

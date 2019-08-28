package git

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"io"
	"os"
	"path/filepath"
)

type Repository struct {
	goGitRepository *git.Repository
}

func NewRepositoryFromFilePath(pathToRepository string) (*Repository, error) {
	// validate given repository path
	pathToRepository, err := filepath.Abs(pathToRepository)
	if err != nil {
		return nil, ErrOpeningRepository{Reasons: []string{
			"getting absolute path",
			err.Error(),
		}}
	}
	if pathToRepository == "" {
		return nil, ErrOpeningRepository{Reasons: []string{
			"given path to repository is blank",
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
	repository, err := git.PlainOpen(pathToRepository)
	if err != nil {
		return nil, ErrOpeningRepository{Reasons: []string{
			err.Error(),
		}}
	}

	return &Repository{
		goGitRepository: repository,
	}, nil
}

func GetDeploymentTag() (string, error) {

	// check if a tag has been provided as a plugin option
	if os.Getenv("PLUGIN_DEPLOYMENT_TAG") != "" {
		return os.Getenv("PLUGIN_DEPLOYMENT_TAG"), nil
	}

	// confirm that given repository path is valid
	repositoryPath := os.Getenv("PLUGIN_DEPLOYMENT_TAG_REPOSITORY_PATH")
	if repositoryPath == "" {
		return "", errors.New("deployment tag repository path is blank")
	}
	f, err := os.Stat(repositoryPath)
	if err != nil {
		return "", errors.New("unable to validate given deployment tag repository path: " + err.Error())
	}
	if !f.IsDir() {
		return "", errors.New(fmt.Sprintf("'%s' is not a directory", repositoryPath))
	}

	// try and open repository at given path
	repository, err := git.PlainOpen(repositoryPath)
	if err != nil {
		return "", errors.New("unable to open repository: " + err.Error())
	}

	// try get the latest tag to use as deployment tag
	tagReferences, err := repository.Tags()
	if err != nil {
		log.Fatal("error getting repository tags: " + err.Error())
	}
	if tagReferences != nil {
		latestTag, err := tagReferences.Next()
		switch err {
		case io.EOF:
			// no tags
		case nil:
			return latestTag.Name().Short(), nil
		default:
			return "", errors.New("error getting latest repository tag ref: " + err.Error())
		}
	}

	// try and get latest commit hash to use as deployment tag
	commitReferences, err := repository.CommitObjects()
	if err != nil {
		log.Fatal("error getting repository commits: " + err.Error())
	}
	commit, err := commitReferences.Next()
	if err != nil {
		return "", errors.New("error getting latest repository commit: " + err.Error())
	}
	return commit.Hash.String(), nil
}

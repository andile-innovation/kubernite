package main

import (
	"errors"
	"fmt"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubeRestclient "k8s.io/client-go/rest"
	"log"
	"os"
	"path/filepath"
)

var repoPath = "/Users/bernardbussy/go/src/github.com/andile-innovation/james"

func main() {

	var config = new(kubeRestclient.Config)
	var err error

	deploymentTag, err := getDeploymentTag()
	if err != nil {
		log.Fatal("error getting deployment tag: " + err.Error())
	}

	fmt.Println("use: " + deploymentTag)

	config.Host = os.Getenv("PLUGIN_KUBERNETES_SERVER")
	config.TLSClientConfig.CAData = []byte(os.Getenv("PLUGIN_KUBERNETES_CERT_DATA"))
	config.TLSClientConfig.CertData = []byte(os.Getenv("PLUGIN_KUBERNETES_CLIENT_CERT_DATA"))
	config.TLSClientConfig.KeyData = []byte(os.Getenv("PLUGIN_KUBERNETES_CLIENT_KEY_DATA"))

	// create the client set
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		pods, err := clientset.CoreV1().Pods("dev").List(v1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster!\n", len(pods.Items))
		break
	}
}

func getDeploymentTag() (string, error) {

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

	r, err := openRepository()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(r)

	// try and open repository at given path
	repository, err := git.PlainOpen(repoPath)
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

func openRepository() (*git.Repository, error) {
	// confirm that given repository path is valid
	repositoryPath := os.Getenv("PLUGIN_DEPLOYMENT_TAG_REPOSITORY_PATH")
	if repositoryPath == "" {
		return nil, errors.New("deployment tag repository path is blank")
	}
	absoluteRepositoryPath, err := filepath.Abs(repositoryPath)
	if err != nil {
		return nil, errors.New("error getting absolute repository path: " + err.Error())
	}
	f, err := os.Stat(absoluteRepositoryPath)
	if err != nil {
		return nil, errors.New("unable to validate given deployment tag repository path: " + err.Error())
	}
	if !f.IsDir() {
		return nil, errors.New(fmt.Sprintf("'%s' is not a directory", absoluteRepositoryPath))
	}

	// get a billy filesystem for repository .git file
	var repositoryFS, dotGitFS billy.Filesystem
	var dotGitFSInfo os.FileInfo
	repositoryFS = osfs.New(absoluteRepositoryPath)
	dotGitFSInfo, err = repositoryFS.Stat(".git")
	if err != nil {
		return nil, errors.New("unable to get .git file/directory info: " + err.Error())
	}
	if !dotGitFSInfo.IsDir() {
		return nil, errors.New("<deployment_tag_repository_path>/.git is not a directory")
	}

	dotGitFS, err = repositoryFS.Chroot(".git")
	if err != nil {
		return nil, errors.New("unable to get .git directory filesystem: " + err.Error())
	}

	if _, err := dotGitFS.Stat(""); err != nil {
		return nil, errors.New("unable to get .git directory filesystem info: " + err.Error())
	}

	return nil, nil
}

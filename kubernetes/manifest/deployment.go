package manifest

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Deployment struct {
	Manifest   `yaml:"-"`
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string `yaml:"name"`
		Namespace   string `yaml:"namespace"`
		Annotations struct {
			KubernetesIOChangeCause string `yaml:"kubernetes.io/change-cause"`
		} `yaml:"annotations"`
	} `yaml:"metadata"`
}

func NewDeploymentFromFile(pathToDeploymentFile string) (*Deployment, error) {
	// validate given path to deployment file
	pathToDeploymentFile, err := filepath.Abs(pathToDeploymentFile)
	if err != nil {
		return nil, ErrParsingManifestFile{Reasons: []string{
			"getting absolute path",
			err.Error(),
		}}
	}
	deploymentFileInfo, err := os.Stat(pathToDeploymentFile)
	if err != nil {
		return nil, ErrParsingManifestFile{Reasons: []string{
			"getting file info",
			err.Error(),
		}}
	}
	if deploymentFileInfo.IsDir() {
		return nil, ErrParsingManifestFile{Reasons: []string{
			fmt.Sprintf("'%s' is a directory", pathToDeploymentFile),
		}}
	}

	// open deployment file
	deploymentFile, err := ioutil.ReadFile(pathToDeploymentFile)
	if err != nil {
		return nil, ErrParsingManifestFile{Reasons: []string{
			"reading deployment file",
			err.Error(),
		}}
	}

	// parse deployment file
	var parsedDeploymentFile Deployment
	if err := yaml.Unmarshal(deploymentFile, &parsedDeploymentFile); err != nil {
		return nil, ErrParsingManifestFile{Reasons: []string{
			"parsing deployment file",
			err.Error(),
		}}
	}

	// set deployment manifest file path
	parsedDeploymentFile.Manifest.PathToFile = pathToDeploymentFile

	return &parsedDeploymentFile, nil
}

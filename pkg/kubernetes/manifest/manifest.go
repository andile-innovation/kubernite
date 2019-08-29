package manifest

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Manifest struct {
	PathToFile  string
	YAMLContent map[string]interface{}
}

func NewManifest(pathToManifestFile string) (*Manifest, error) {
	// validate given path to deployment file
	pathToManifestFile, err := filepath.Abs(pathToManifestFile)
	if err != nil {
		return nil, ErrParsingManifestFile{Reasons: []string{
			"getting absolute path",
			err.Error(),
		}}
	}
	deploymentFileInfo, err := os.Stat(pathToManifestFile)
	if err != nil {
		return nil, ErrParsingManifestFile{Reasons: []string{
			"getting file info",
			err.Error(),
		}}
	}
	if deploymentFileInfo.IsDir() {
		return nil, ErrParsingManifestFile{Reasons: []string{
			fmt.Sprintf("'%s' is a directory", pathToManifestFile),
		}}
	}

	// open deployment file
	manifestFile, err := ioutil.ReadFile(pathToManifestFile)
	if err != nil {
		return nil, ErrParsingManifestFile{Reasons: []string{
			"reading manifest file",
			err.Error(),
		}}
	}

	// instantiate new manifest object
	newManifest := new(Manifest)
	newManifest.YAMLContent = make(map[string]interface{})

	// parse manifest file
	if err := yaml.Unmarshal(manifestFile, &(newManifest.YAMLContent)); err != nil {
		return nil, ErrParsingManifestFile{Reasons: []string{
			"parsing deployment file",
			err.Error(),
		}}
	}

	// set deployment manifest file path
	newManifest.PathToFile = pathToManifestFile

	return newManifest, nil
}

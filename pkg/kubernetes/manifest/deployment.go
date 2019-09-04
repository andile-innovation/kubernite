package manifest

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	v1 "k8s.io/api/apps/v1"
	k8sYamlUtil "k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
)

/*
Deployment is a convenience wrapper the manifest object type that represents a deployment file
*/
type Deployment struct {
	*v1.Deployment
	PathToFile string
}

/*
NewDeploymentFromFile creates a new deployment file wrapper from a file located at a
given file path.
*/
func NewDeploymentFromFile(pathToDeploymentFile string) (*Deployment, error) {
	// validate given path to deployment file
	pathToDeploymentFile, err := filepath.Abs(pathToDeploymentFile)
	if err != nil {
		return nil, ErrInvalidFilePath{Reasons: []string{
			"could not convert to absolute path",
			err.Error(),
		}}
	}
	deploymentFileInfo, err := os.Stat(pathToDeploymentFile)
	if err != nil {
		return nil, ErrInvalidFilePath{Reasons: []string{
			"could not get file info at path",
			err.Error(),
		}}
	}
	if deploymentFileInfo.IsDir() {
		return nil, ErrInvalidFilePath{Reasons: []string{
			fmt.Sprintf("'%s' is a directory", pathToDeploymentFile),
		}}
	}

	// instantiate a new deployment file and set path
	newDeployment := new(Deployment)
	newDeployment.Deployment = new(v1.Deployment)
	newDeployment.PathToFile = pathToDeploymentFile

	// open a file reader for the deployment file
	deploymentFileReader, err := os.Open(pathToDeploymentFile)
	if err != nil {
		return nil, ErrUnexpected{Reasons: []string{
			"opening deployment file",
			err.Error(),
		}}
	}
	defer func() {
		if err := deploymentFileReader.Close(); err != nil {
			log.Error("closing deployment file: " + err.Error())
		}
	}()

	// decode the deployment yaml file
	if err := k8sYamlUtil.NewYAMLOrJSONDecoder(deploymentFileReader, 512).Decode(&newDeployment.Deployment); err != nil {
		return nil, ErrUnexpected{Reasons: []string{
			"decoding deployment file",
			err.Error(),
		}}
	}

	return newDeployment, nil
}

func (d *Deployment) UpdateAnnotations(key, value string) error {
	if d.Annotations == nil {
		d.Annotations = make(map[string]string)
	}
	d.Annotations[key] = value
	return nil
}

func (d *Deployment) UpdatePodTemplateAnnotations(key, value string) error {
	if d.Spec.Template.Annotations == nil {
		d.Spec.Template.Annotations = make(map[string]string)
	}
	d.Spec.Template.Annotations[key] = value
	return nil
}

func (d *Deployment) UpdateImageTag(deploymentImageName, latestTag string) error {
	// validation
	if len(d.Spec.Template.Spec.Containers) == 0 {
		return ErrDeploymentManifestInvalid{
			Reasons: []string{
				"no images in pod spec",
			},
		}
	}

	if len(d.Spec.Template.Spec.Containers) == 1 {
		containerImage := d.Spec.Template.Spec.Containers[0].Image
		if deploymentImageName == "" {
			deploymentImageName = containerImage[:strings.IndexByte(containerImage, ':')]
		} else {
			if deploymentImageName != containerImage[:strings.IndexByte(containerImage, ':')] {
				return ErrSuppliedImageNameNotInConfigFile{}
			}
		}
		d.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s",deploymentImageName,latestTag)
		return nil
	}

	if deploymentImageName == "" {
		return ErrImageNotSpecified{}
	}

	for i, c := range d.Spec.Template.Spec.Containers {
		if c.Image[:strings.IndexByte(c.Image, ':')] == deploymentImageName {
			d.Spec.Template.Spec.Containers[i].Image = fmt.Sprintf("%s:%s",deploymentImageName,latestTag)
			return nil
		}
	}
	return ErrSuppliedImageNameNotInConfigFile{}
}

/*
WriteToYAML writes the manifest file to disk at it's original filepath
*/
func (d *Deployment) WriteToYAML() error {
	return d.WriteToYAMLAtPath(d.PathToFile)
}

/*
WriteAtPath writes the manifest file to disk at given file path
*/
func (d *Deployment) WriteToYAMLAtPath(pathToWriteManifestFile string) error {
	// confirm that file path has correct extension
	if !(strings.HasSuffix(pathToWriteManifestFile, ".yaml") || strings.HasSuffix(pathToWriteManifestFile, ".yml")) {
		return ErrInvalidFilePath{Reasons: []string{
			fmt.Sprintf("'%s' does not end in .yaml or .yml", pathToWriteManifestFile),
		}}
	}

	// marshal deployment object to json
	jsonData, err := json.Marshal(d.Deployment)
	if err != nil {
		return ErrUnexpected{Reasons: []string{
			"marshalling to json",
			err.Error(),
		}}
	}

	// convert json data to yaml data
	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return ErrUnexpected{Reasons: []string{
			"converting json to yaml",
			err.Error(),
		}}
	}

	// write to file
	if err := ioutil.WriteFile(pathToWriteManifestFile, yamlData, 0644); err != nil {
		return ErrUnexpected{Reasons: []string{
			"writing to file",
			err.Error(),
		}}
	}

	return nil
}

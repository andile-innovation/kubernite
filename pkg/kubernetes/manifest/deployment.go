package manifest

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	k8sYamlUtil "k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path/filepath"
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

/*
WriteAtPath writes the manifest file to disk at given file path
*/
func (d *Deployment) WriteAtPath(pathToWriteManifestFile string) error {
	deploymentFileBytes, err := d.Marshal()
	if err != nil {
		return ErrUnexpected{Reasons: []string{
			"marshalling deployment object",
			err.Error(),
		}}
	}

	fmt.Println(string(deploymentFileBytes))

	return nil
}

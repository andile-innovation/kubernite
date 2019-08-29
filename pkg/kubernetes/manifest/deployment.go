package manifest

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Deployment struct {
	*Manifest
}

func (d *Deployment) UpdateChangeCauseAnnotation(newValue string) error {
	return nil
}

func NewDeploymentFromFile(pathToDeploymentFile string) (*Deployment, error) {
	newManifest, err := NewManifest(pathToDeploymentFile)
	if err != nil {
		return nil, err
	}

	// instantiate new deployment file using manifest
	newDeployment := new(Deployment)
	newDeployment.Manifest = newManifest

	// confirm that manifest kind is Deployment
	if newDeployment.Kind != DeploymentKind {
		return nil, ErrInvalidManifestKind{
			Expected: DeploymentKind,
			Actual:   newDeployment.Kind,
		}
	}

	output, err := yaml.Marshal(newDeployment.YAMLContent)
	if err != nil {
		log.Fatal("error marshalling!")
	}
	err = ioutil.WriteFile("output.yaml", output, 0644)
	if err != nil {
		log.Fatal("error writing out!!")
	}

	return newDeployment, nil
}

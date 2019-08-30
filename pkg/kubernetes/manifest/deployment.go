package manifest

import (
	"fmt"
	k8sYamlUtil "k8s.io/apimachinery/pkg/util/yaml"
)

/*
Deployment is a convenience wrapper the manifest object type that represents a deployment file
*/
type Deployment struct {
	*Manifest
}

/*
NewDeploymentFromFile creates a new deployment file wrapper from a file located at a
given file path.
*/
func NewDeploymentFromFile(pathToDeploymentFile string) (*Deployment, error) {
	newManifest, err := NewManifestFromFile(pathToDeploymentFile)
	if err != nil {
		return nil, err
	}

	// instantiate new deployment file using manifest
	newDeployment := new(Deployment)
	newDeployment.Manifest = newManifest

	// confirm that manifest kind is Deployment
	if newDeployment.Kind != DeploymentKind {
		return nil, ErrDeploymentManifestInvalid{
			Reasons: []string{
				fmt.Sprintf(
					"incorrect manifest kind '%s' - should be '%s' ",
					newDeployment.Kind,
					DeploymentKind,
				),
			},
		}
	}

	return newDeployment, nil
}

func (d *Deployment) UpdateAnnotations(key, value string) error {
	// find annotations section
	annotationsSectionMap, err := d.GetObjectMap("metadata.annotations")
	if err != nil {
		switch err.(type) {
		case ErrKeyNotFoundInObject:
			// if annotations is not found in metadata, get metadata section
			metadataSectionMap, err := d.GetObjectMap("metadata")
			if err != nil {
				return err
			}
			// add annotation section to it
			(*metadataSectionMap)["annotations"] = make(map[interface{}]interface{})
			// find it again
			annotationsSectionMap, err = d.GetObjectMap("metadata.annotations")
			if err != nil {
				return ErrUnexpected{Reasons: []string{
					"unable to add annotations section to deployment file",
					err.Error(),
				}}
			}
		default:
			return err
		}
	}

	// update annotation
	(*annotationsSectionMap)[key] = value

	return nil
}

func (d *Deployment) UpdatePodTemplateAnnotations(key, value string) error {
	// find annotations section
	annotationsSectionMap, err := d.GetObjectMap("spec.template.metadata.annotations")
	if err != nil {
		switch err.(type) {
		case ErrKeyNotFoundInObject:
			// if annotations is not found in metadata, get metadata section
			metadataSectionMap, err := d.GetObjectMap("spec.template.metadata")
			if err != nil {
				return err
			}
			// add annotation section to it
			(*metadataSectionMap)["annotations"] = make(map[interface{}]interface{})
			// find it again
			annotationsSectionMap, err = d.GetObjectMap("spec.template.metadata.annotations")
			if err != nil {
				return ErrUnexpected{Reasons: []string{
					"unable to add annotations section to pod template",
					err.Error(),
				}}
			}
		default:
			return err
		}
	}

	// update annotation
	(*annotationsSectionMap)[key] = value

	return nil
}

func (d *Deployment) ToJSON() ([]byte, error) {
	deploymentFileContents, err := d.GetDeploymentFileContents()
	if err != nil {
		return nil, ErrUnexpected{Reasons: []string{
			"getting deployment file contents",
		}}
	}

	jsonContent, err := k8sYamlUtil.ToJSON(deploymentFileContents)
	if err != nil {
		return nil, ErrUnexpected{Reasons: []string{
			"converting to json",
		}}
	}

	return jsonContent, nil
}

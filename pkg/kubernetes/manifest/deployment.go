package manifest

import "fmt"

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

	fmt.Println(annotationsSectionMap)

	(*annotationsSectionMap)[key] = value

	return nil
}

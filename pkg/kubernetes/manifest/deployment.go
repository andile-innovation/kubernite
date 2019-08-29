package manifest

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
		return nil, ErrInvalidManifestKind{
			Expected: DeploymentKind,
			Actual:   newDeployment.Kind,
		}
	}

	return newDeployment, nil
}

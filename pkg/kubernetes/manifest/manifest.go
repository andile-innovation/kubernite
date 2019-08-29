package manifest

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Kind string

func (k Kind) String() string {
	return string(k)
}

const DeploymentKind Kind = "Deployment"

type Manifest struct {
	PathToFile  string
	Kind        Kind
	YAMLContent map[interface{}]interface{}
}

func NewManifestFromFile(pathToManifestFile string) (*Manifest, error) {
	// validate given path to deployment file
	pathToManifestFile, err := filepath.Abs(pathToManifestFile)
	if err != nil {
		return nil, ErrInvalidFilePath{Reasons: []string{
			"could not convert to absolute path",
			err.Error(),
		}}
	}
	deploymentFileInfo, err := os.Stat(pathToManifestFile)
	if err != nil {
		return nil, ErrInvalidFilePath{Reasons: []string{
			"could not get file info at path",
			err.Error(),
		}}
	}
	if deploymentFileInfo.IsDir() {
		return nil, ErrInvalidFilePath{Reasons: []string{
			fmt.Sprintf("'%s' is a directory", pathToManifestFile),
		}}
	}

	// open deployment file
	manifestFile, err := ioutil.ReadFile(pathToManifestFile)
	if err != nil {
		return nil, ErrUnexpected{Reasons: []string{
			"reading manifest file",
			err.Error(),
		}}
	}

	// instantiate new manifest object
	newManifest := new(Manifest)
	newManifest.YAMLContent = make(map[interface{}]interface{})

	// parse manifest file
	if err := yaml.UnmarshalStrict(manifestFile, &(newManifest.YAMLContent)); err != nil {
		return nil, ErrUnexpected{Reasons: []string{
			"parsing deployment file",
			err.Error(),
		}}
	}

	// set deployment manifest file path
	newManifest.PathToFile = pathToManifestFile

	// set manifest kind and populate top level objects
	for key := range newManifest.YAMLContent {
		object, found := newManifest.YAMLContent[key]
		if !found {
			return nil, ErrUnexpected{Reasons: []string{
				"no object found at key",
			}}
		}

		// look for the kind key
		if key == "kind" {
			kind, ok := object.(string)
			if !ok {
				return nil, ErrUnexpected{Reasons: []string{
					"inferring type of manifest file kind field",
				}}
			}
			newManifest.Kind = Kind(kind)
		}
	}

	// confirm that manifest kind is valid
	switch newManifest.Kind {
	case DeploymentKind:
	default:
		return nil, ErrDeploymentManifestInvalid{
			Reasons: []string{
				fmt.Sprintf("invalid/unsupported manifest kind: '%s'", newManifest.Kind),
			},
		}
	}

	return newManifest, nil
}

/*
Write writes the manifest file to disk at its deployment.Manifest.PathToFile
*/
func (m *Manifest) Write() error {
	return m.WriteAtPath(m.PathToFile)
}

/*
WriteAtPath writes the manifest file to disk at given file path
*/
func (m *Manifest) WriteAtPath(pathToWriteManifestFile string) error {
	// marshall manifest
	marshalledYAML, err := yaml.Marshal(m.YAMLContent)
	if err != nil {
		return ErrUnexpected{Reasons: []string{
			"marshalling yaml content",
			err.Error(),
		}}
	}

	// confirm that given path is valid (i.e. ends with .yaml or .yml)
	if !(strings.HasSuffix(m.PathToFile, ".yaml") || strings.HasSuffix(m.PathToFile, ".yml")) {
		return ErrInvalidFilePath{Reasons: []string{
			fmt.Sprintf("'%s' does not end with .yaml or .yml", m.PathToFile),
		}}
	}

	// write to disk
	if err := ioutil.WriteFile(pathToWriteManifestFile, marshalledYAML, 0644); err != nil {
		return ErrUnexpected{Reasons: []string{
			"writing file to disk",
			err.Error(),
		}}
	}

	return nil
}

func (m *Manifest) GetObjectMap(accessor string) (*map[interface{}]interface{}, error) {
	if m.YAMLContent == nil {
		return nil, ErrUnexpected{Reasons: []string{
			"manifest YAML Content is nil",
		}}
	}
	return GetObjectMapAtKey(&m.YAMLContent, accessor)
}

func GetObjectMapAtKey(objectToSearch *map[interface{}]interface{}, accessor string) (*map[interface{}]interface{}, error) {
	// split accessor path
	accessorKeys := strings.Split(accessor, ".")
	if len(accessorKeys) == 0 {
		return nil, ErrUnexpected{Reasons: []string{
			"zero length accessor path",
		}}
	}

	// get the key to find and build remaining accessor path
	keyToFind := accessorKeys[0]
	var remainingAccessorPath string
	if len(accessorKeys) > 1 {
		remainingAccessorPath = strings.Join(accessorKeys[1:], ".")
	}

	// find object addressed by first accessor
	object, found := (*objectToSearch)[keyToFind]
	if !found {
		return nil, ErrKeyNotFoundInObject{
			Key:    keyToFind,
			Object: *objectToSearch,
		}
	}

	// infer object type
	objectMap, ok := object.(map[interface{}]interface{})
	if !ok {
		return nil, ErrUnexpected{Reasons: []string{
			"inferring type of object map",
		}}
	}

	// check if we are at the end of the accessor path
	if len(remainingAccessorPath) == 0 {
		return &objectMap, nil
	}

	// if we are not call this function again with the remaining accessors
	return GetObjectMapAtKey(&objectMap, remainingAccessorPath)
}

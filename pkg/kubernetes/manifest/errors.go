package manifest

import (
	"fmt"
	"strings"
)

type ErrInvalidFilePath struct {
	Reasons []string
}

func (e ErrInvalidFilePath) Error() string {
	return "invalid manifest file path: " + strings.Join(e.Reasons, ", ")
}

type ErrUnexpected struct {
	Reasons []string
}

func (e ErrUnexpected) Error() string {
	return "unexpected manifest file error: " + strings.Join(e.Reasons, ", ")
}

type ErrManifestInvalid struct {
	Reasons []string
}

func (e ErrManifestInvalid) Error() string {
	return "manifest invalid: " + strings.Join(e.Reasons, ", ")
}

type ErrDeploymentManifestInvalid struct {
	Reasons []string
}

func (e ErrDeploymentManifestInvalid) Error() string {
	return "deployment manifest invalid: " + strings.Join(e.Reasons, ", ")
}

type ErrInvalidAccessorPath struct {
	AccessorPath string
	Object       interface{}
}

func (e ErrInvalidAccessorPath) Error() string {
	return fmt.Sprintf("invalid accessor path '%s' for object %v", e.AccessorPath, e.Object)
}

type ErrKeyNotFoundInObject struct {
	Key    string
	Object interface{}
}

func (e ErrKeyNotFoundInObject) Error() string {
	return fmt.Sprintf("key '%s' not found in object %v", e.Key, e.Object)
}

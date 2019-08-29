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

type ErrInvalidManifestKind struct {
	Expected Kind
	Actual   Kind
}

func (e ErrInvalidManifestKind) Error() string {
	return fmt.Sprintf(
		"invalid manifest kind - expected '%s' vs actual '%s'",
		e.Expected,
		e.Actual,
	)
}

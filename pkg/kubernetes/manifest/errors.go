package manifest

import (
	"fmt"
	"strings"
)

type ErrParsingManifestFile struct {
	Reasons []string
}

func (e ErrParsingManifestFile) Error() string {
	return "error parsing manifest file: " + strings.Join(e.Reasons, ", ")
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

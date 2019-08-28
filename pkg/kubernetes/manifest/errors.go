package manifest

import "strings"

type ErrParsingManifestFile struct {
	Reasons []string
}

func (e ErrParsingManifestFile) Error() string {
	return "error parsing manifest file: " + strings.Join(e.Reasons, ", ")
}

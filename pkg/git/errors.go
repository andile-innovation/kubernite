package git

import "strings"

type ErrOpeningRepository struct {
	Reasons []string
}

func (e ErrOpeningRepository) Error() string {
	return "error opening repository: " + strings.Join(e.Reasons, ", ")
}

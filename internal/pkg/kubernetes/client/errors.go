package client

import "strings"

type ErrCreatingClientSet struct {
	Reasons []string
}

func (e ErrCreatingClientSet) Error() string {
	return "error creating client set: " + strings.Join(e.Reasons, ", ")
}

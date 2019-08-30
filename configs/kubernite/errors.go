package kubernite

import "strings"

type ErrPackageInitialisation struct {
	Reasons []string
}

func (e ErrPackageInitialisation) Error() string {
	return "package initialisation error: " + strings.Join(e.Reasons, ", ")
}

type ErrInvalidConfig struct {
	Reasons []string
}

func (e ErrInvalidConfig) Error() string {
	return "invalid config:\n" + strings.Join(e.Reasons, ", ")
}

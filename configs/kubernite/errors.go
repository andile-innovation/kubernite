package kubernite

import "strings"

type ErrInvalidConfig struct {
	Reasons []string
}

func (e ErrInvalidConfig) Error() string {
	return "invalid config:\n" + strings.Join(e.Reasons, ",")
}

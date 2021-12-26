package core

import "github.com/pkg/errors"

type ErrNotFound struct{}

func (n ErrNotFound) Error() string {
	return "not found"
}

func IsErrNotFound(err error) bool {
	switch errors.Cause(err).(type) {
	case *ErrNotFound:
		return true
	default:
		return false
	}
}

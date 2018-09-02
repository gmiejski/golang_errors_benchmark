package src

import "github.com/pkg/errors"

type ErrorWrapper = func(err error, msg string) error

var errorsWrapperUsed = withWrapper

var withWrapper = func(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

var wrapWrapper = func(err error, msg string) error {
	return errors.Wrap(err, msg)
}


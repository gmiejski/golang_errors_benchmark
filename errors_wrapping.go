package errors_benchmarks

import "github.com/pkg/errors"

type ErrorWrapper = func(err error, msg string) error

var defaultErrWrapper = withWrapper

var withWrapper = func(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

var wrapWrapper = func(err error, msg string) error {
	return errors.Wrap(err, msg)
}


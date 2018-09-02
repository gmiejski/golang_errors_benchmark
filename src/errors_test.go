package src

import (
	"testing"
	"github.com/pkg/errors"
)





func BenchmarkErrors(b *testing.B) {

	err := errors.New("some error")
	for i := 0; i < b.N; i++ {
		err = errors.WithMessage(err, "a")
	}
}

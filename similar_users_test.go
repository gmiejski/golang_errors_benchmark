package errors_benchmarks

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"
	"fmt"
	"strings"
	stdErrors "errors"
)

var expectedErrorMsg = "cannot find similar user to user -1: repository cannot find similar users: error calling db driver"
var expectedCauseMsg = "error calling db driver"

var serviceWithWrappers = SimilarUserService{errWrapper: withWrapper, repo: SimilarUserRepo{errWrapper: withWrapper, DBDriver: DBDriver{}}}
var serviceWrapWrappers = SimilarUserService{errWrapper: wrapWrapper, repo: SimilarUserRepo{errWrapper: wrapWrapper, DBDriver: DBDriver{}}}
var serviceMixedWrappers = SimilarUserService{errWrapper: withWrapper, repo: SimilarUserRepo{errWrapper: wrapWrapper, DBDriver: DBDriver{}}}

var errorsStrategiesTestCases = []struct {
	name    string
	service SimilarUserService
}{
	{"all-with wrapper", serviceWithWrappers},
	{"all-wrap wrapper", serviceWrapWrappers},
	{"mixed", serviceMixedWrappers},
}

func TestErrorAndCausesAreTheSame(t *testing.T) {
	for _, test := range errorsStrategiesTestCases {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.service.Find(1)
			assert.Error(t, err)
			assert.EqualValues(t, expectedErrorMsg, err.Error())
			assert.EqualValues(t, expectedCauseMsg, errors.Cause(err).Error())
		})
	}
}

func TestStackTraceContainsAllInfoNeeded(t *testing.T) {
	tests := []struct {
		name             string
		service          SimilarUserService
		stacktracesCount int
	}{
		{"all-with wrapper", serviceWithWrappers, 0},
		{"all-wrap wrapper", serviceWrapWrappers, 2},
		{"mixed", serviceMixedWrappers, 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.service.Find(1)

			stacktrace := fmt.Sprintf("%+v", err)
			assert.EqualValues(t, test.stacktracesCount, strings.Count(stacktrace, "github.com/brainly/errors-benchmarks.(*SimilarUserService).Find"))

			fmt.Println(stacktrace)
		})
	}
}

func BenchmarkSimilarUser_Find(b *testing.B) {
	for _, bm := range errorsStrategiesTestCases {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bm.service.Find(1)
			}
		})
	}
}

func BenchmarkErrorsWrappersWithoutContext(b *testing.B) {
	benchmarks := []struct {
		name    string
		wrapper func(err error, mgs string) error
	}{
		{"errors.WithMessage", errors.WithMessage},
		{"errors.Wrap", errors.Wrap},
	}
	for _, bm := range benchmarks {
		var err error
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err = bm.wrapper(stdErrors.New("err"), "err2")
			}
		})
	}
}

func BenchmarkDoubleErrorsWrappersWithoutContext(b *testing.B) {
	benchmarks := []struct {
		name    string
		wrapper func(err error, mgs string) error
	}{
		{"errors.WithMessage", errors.WithMessage},
		{"errors.Wrap", errors.Wrap},
	}
	for _, bm := range benchmarks {
		var err error
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err = bm.wrapper(stdErrors.New("err"), "err2")
			}
		})
	}
}

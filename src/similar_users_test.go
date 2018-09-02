package src

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/pkg/errors"
	"fmt"
)

var expectedStackTrace = `error calling db driver
repository cannot find similar users
github.com/brainly/errors-benchmarks/src.glob..func2
	/Users/grzegorzmiejski/workspaces/go/src/github.com/brainly/errors-benchmarks/src/errors_wrapping.go:14
github.com/brainly/errors-benchmarks/src.(*SimilarUserRepo).Find
	/Users/grzegorzmiejski/workspaces/go/src/github.com/brainly/errors-benchmarks/src/similar_users.go:31
github.com/brainly/errors-benchmarks/src.(*SimilarUserService).Find
	/Users/grzegorzmiejski/workspaces/go/src/github.com/brainly/errors-benchmarks/src/similar_users.go:21
github.com/brainly/errors-benchmarks/src.TestSameErrorsMessages
	/Users/grzegorzmiejski/workspaces/go/src/github.com/brainly/errors-benchmarks/src/similar_users_test.go:27
testing.tRunner
	/usr/local/Cellar/go/1.10.3/libexec/src/testing/testing.go:777
runtime.goexit
	/usr/local/Cellar/go/1.10.3/libexec/src/runtime/asm_amd64.s:2361
cannot find similar user to user -1
github.com/brainly/errors-benchmarks/src.glob..func2
	/Users/grzegorzmiejski/workspaces/go/src/github.com/brainly/errors-benchmarks/src/errors_wrapping.go:14
github.com/brainly/errors-benchmarks/src.(*SimilarUserService).Find
	/Users/grzegorzmiejski/workspaces/go/src/github.com/brainly/errors-benchmarks/src/similar_users.go:22
github.com/brainly/errors-benchmarks/src.TestSameErrorsMessages
	/Users/grzegorzmiejski/workspaces/go/src/github.com/brainly/errors-benchmarks/src/similar_users_test.go:27
testing.tRunner
	/usr/local/Cellar/go/1.10.3/libexec/src/testing/testing.go:777
runtime.goexit
	/usr/local/Cellar/go/1.10.3/libexec/src/runtime/asm_amd64.s:2361`


func TestSameErrorsMessages(t *testing.T) {
	expectedErrorMsg := "cannot find similar user to user -1: repository cannot find similar users: error calling db driver"
	expectedCauseMsg := "error calling db driver"
	service := SimilarUserService{repo: SimilarUserRepo{DBDriver: DBDriver{}}}
	errorsWrapperUsed = withWrapper

	_, errWithWrapper := service.Find(1)
	
	assert.Error(t, errWithWrapper)
	assert.EqualValues(t, expectedErrorMsg, errWithWrapper.Error())
	assert.EqualValues(t, expectedCauseMsg, errors.Cause(errWithWrapper).Error())

	errorsWrapperUsed = wrapWrapper
	
	_, errWrapWrapper := service.Find(1)

	assert.Error(t, errWrapWrapper)
	assert.EqualValues(t, expectedErrorMsg, errWrapWrapper.Error())
	assert.EqualValues(t, expectedCauseMsg, errors.Cause(errWrapWrapper).Error())

	//st1 := fmt.Sprintf("%+v", errWrapWrapper)
	st2 := fmt.Sprintf("%+v", errWithWrapper)

	//assert.EqualValues(t, expectedStackTrace, st1)
	assert.EqualValues(t, `error calling db driver
repository cannot find similar users
cannot find similar user to user -1`, st2)
}

func TestErrorGettingSimilarUsers(t *testing.T) {
	service := SimilarUserService{repo: SimilarUserRepo{DBDriver: DBDriver{}}}
	errorsWrapperUsed = wrapWrapper

	user, err := service.Find(1)

	assert.EqualValues(t, UserNotFound, user)
	assert.Error(t, err)
	println(errors.Cause(err).Error())
	fmt.Printf("%+v", err)
}

func BenchmarkSimilarUser(b *testing.B) {
	service := SimilarUserService{repo: SimilarUserRepo{DBDriver: DBDriver{}}}
	errorsWrapperUsed = wrapWrapper
	for i := 0; i < b.N; i++ {
		service.Find(1)
	}
}

func BenchmarkSimilarUser_Find(b *testing.B) {
	benchmarks := []struct {
		name string
		errFunc ErrorWrapper
	}{
		{"with",withWrapper},
		{"wrap",wrapWrapper},
	}
	for _, bm := range benchmarks {
		errorsWrapperUsed = bm.errFunc
		service := SimilarUserService{repo: SimilarUserRepo{DBDriver: DBDriver{}}}

		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				service.Find(1)
			}
		})
	}
}

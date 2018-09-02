package src

import (
	standardErrors "errors"
	"fmt"
)

type User = int

var UserNotFound = -1

type SimilarUser interface {
	Find(user User) (User, error)
}

type SimilarUserService struct {
	repo       SimilarUserRepo
	errWrapper ErrorWrapper
}

func (s *SimilarUserService) Find(user User) (User, error) {
	user, err := s.repo.Find(user)
	return user, errorExtendOrDefault(err, fmt.Sprintf("cannot find similar user to user %d", user), s.errWrapper)
}

type SimilarUserRepo struct {
	DBDriver   DBDriver
	errWrapper ErrorWrapper
}

func (s *SimilarUserRepo) Find(user User) (User, error) {
	user, err := s.DBDriver.Find(user)
	return user, errorExtendOrDefault(err, "repository cannot find similar users", s.errWrapper)
}

type DBDriver struct {
}

// call some external code - probably will use standard errors package
func (s *DBDriver) Find(user User) (User, error) {
	return UserNotFound, standardErrors.New("error calling db driver")
}

func errorExtendOrDefault(err error, msg string, wrapper ErrorWrapper) error {
	if wrapper != nil {
		return wrapper(err, msg)
	}
	return defaultErrWrapper(err, msg)
}

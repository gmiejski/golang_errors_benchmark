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
	repo SimilarUserRepo
}

func (s *SimilarUserService) Find(user User) (User, error) {
	user, err := s.repo.Find(user)
	return user, errorsWrapperUsed(err, fmt.Sprintf("cannot find similar user to user %d", user))
}

type SimilarUserRepo struct {
	DBDriver DBDriver
}

func (s *SimilarUserRepo) Find(user User) (User, error) { // TODO change to simple userNotFound in DB
	user, err := s.DBDriver.Find(user)
	return user, errorsWrapperUsed(err, "repository cannot find similar users")
}

type DBDriver struct {
}

// call some external code - probably will use standard errors package
func (s *DBDriver) Find(user User) (User, error) {
	return UserNotFound, standardErrors.New("error calling db driver")
}

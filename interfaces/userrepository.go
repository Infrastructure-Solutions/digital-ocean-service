package interfaces

import (
	"errors"

	"github.com/Tinker-Ware/digital-ocean-service/domain"
)

type FakeUserRepo struct {
}

func NewUserRepo() *FakeUserRepo {
	return &FakeUserRepo{}
}

var counter = 0
var users = []domain.User{}

func (repo FakeUserRepo) Store(user domain.User) error {
	user.ID = int64(counter)
	counter++
	users = append(users, user)
	return nil
}

func (repo FakeUserRepo) RetrieveByID(id int64) (*domain.User, error) {

	for _, user := range users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("User not found")
}

package daomock

import (
	"github.com/stretchr/testify/mock"
	"go-tour/internal/models"
)

type UserDaoMock struct {
	mock.Mock
}

func (u *UserDaoMock) FindByID(id uint) (*models.User, error) {
	args := u.Called(id)

	return args.Get(0).(*models.User), args.Error(1)
}

func (u *UserDaoMock) FindByEmail(email string) (*models.User, error) {
	args := u.Called(email)

	return args.Get(0).(*models.User), args.Error(1)
}

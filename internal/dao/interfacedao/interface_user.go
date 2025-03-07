package interfacedao

import "go-tour/internal/models"

type UserDaoInterface interface {
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
}

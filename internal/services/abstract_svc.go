package services

import (
	"context"
	"go-tour/common"
	"go-tour/internal/dao/interfacedao"
	"go-tour/internal/models"
	"go-tour/internal/must"
	"go-tour/internal/serializers"
)

type AbstractService struct {
}

func (e *AbstractService) userFromContext(c context.Context, userDao interfacedao.UserDaoInterface) (*models.User, error) {
	authen := c.Value(common.CustomerKey)
	if authen == nil {
		return nil, must.ErrInvalidCredentials
	}

	userIDVal := authen.(*serializers.UserInfo)

	user, err := userDao.FindByID(userIDVal.ID)
	if err != nil || user == nil {
		return nil, must.ErrInvalidCredentials
	}

	return user, nil
}

package services

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-tour/config"
	"go-tour/gen"
	"go-tour/internal/dao/interfacedao"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type userServiceInterface interface {
	gen.UserServiceServer
}

var _ userServiceInterface = (*UserService)(nil)

type UserService struct {
	AbstractService
	logger  *zap.Logger
	cfg     *config.Config
	db      *gorm.DB
	userDao interfacedao.UserDaoInterface
}

func NewUserService(logger *zap.Logger, cfg *config.Config, db *gorm.DB, userDao interfacedao.UserDaoInterface) *UserService {
	return &UserService{logger: logger, cfg: cfg, db: db, userDao: userDao}
}

func (u *UserService) RegisterGrpcServer(
	s *grpc.Server,
) {
	gen.RegisterUserServiceServer(s, u)
}

func (u *UserService) RegisterHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	conn *grpc.ClientConn,
) error {
	if err := gen.RegisterUserServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	return nil
}

func (u *UserService) Profile(ctx context.Context, in *gen.EmptyRequest) (*gen.UserInfoResponse, error) {
	user, err := u.userFromContext(ctx, u.userDao)
	if err != nil {
		return nil, err
	}

	resp := &gen.UserInfoResponse{
		ID:        int64(user.ID),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return resp, nil
}

package services

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go-tour/config"
	"go-tour/gen"
	"go-tour/internal/dao/interfacedao"
	"go-tour/internal/models"
	"go-tour/internal/must"
	"go-tour/internal/serializers"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

var _ gen.AuthenticationServiceServer = &AuthenticationService{}

type AuthenticationService struct {
	AbstractService
	logger  *zap.Logger
	cfg     *config.Config
	userDao interfacedao.UserDaoInterface
}

func NewAuthenticationService(logger *zap.Logger, cfg *config.Config, userDao interfacedao.UserDaoInterface) *AuthenticationService {
	return &AuthenticationService{logger: logger, cfg: cfg, userDao: userDao}
}

func (a *AuthenticationService) RegisterGrpcServer(
	s *grpc.Server,
) {
	gen.RegisterAuthenticationServiceServer(s, a)
}

func (a *AuthenticationService) RegisterHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	conn *grpc.ClientConn,
) error {
	if err := gen.RegisterAuthenticationServiceHandler(ctx, mux, conn); err != nil {
		return err
	}

	return nil
}

func (a *AuthenticationService) AuthFuncOverride(
	ctx context.Context,
	fullMethodName string,
) (context.Context, error) {
	return ctx, nil
}

func (a *AuthenticationService) SayHello(ctx context.Context, in *gen.HelloRequest) (*gen.HelloReply, error) {
	return &gen.HelloReply{Message: in.Name + " world"}, nil
}

func (e *AuthenticationService) Auth(ctx context.Context, in *gen.LoginRequest) (*gen.LoginResponse, error) {
	user, err := e.authenticatorByEmailPassword(in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	data := &serializers.UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	expire := time.Now().Add(60 * time.Minute)
	accessToken, err := must.CreateNewWithClaims(data, e.cfg.AuthenticationSecretKey, expire)

	return &gen.LoginResponse{
		Data: &gen.LoginResponse_Data{
			AccessToken:  accessToken,
			RefreshToken: "",
			ExpiredIn:    fmt.Sprintf("%d", expire.Unix()),
		},
	}, nil
}

func (u *AuthenticationService) isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func (u *AuthenticationService) authenticatorByEmailPassword(email, password string) (*models.User, error) {
	if !u.isEmailValid(email) {
		return nil, must.ErrInvalidEmail
	}

	user, _ := u.userDao.FindByEmail(email)
	if user != nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return nil, must.ErrInvalidPassword
		}

		return user, nil
	}

	return nil, must.ErrEmailNotExists
}

package service

import (
	"context"
	"ecom/internal/user/dto"
	"ecom/internal/user/entity"
	"ecom/internal/user/repository"
	"ecom/pkg/jwt"
	"ecom/pkg/utils"
	"errors"

	"github.com/quangdangfit/gocommon/logger"
	"github.com/quangdangfit/gocommon/validation"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	Register(ctx context.Context, req *dto.RegisterReq) (*entity.User, error)
}

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(
	validator validation.Validation,
	repo *repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginReq) (*entity.User, string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Errorf("Login.GetUserByEmail fail, email: %s, error: %s", req.Email, err)
		return nil, "", "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", "", errors.New("Invalid credentials...")
	}
	tokenData := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	}
	accessToken := jwt.GenerateAccessToken(tokenData)
	refreshToken := jwt.GenerateRefreshToken(tokenData)

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterReq) (*entity.User, error) {
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		logger.Errorf("User already exits for email: %s, error: %s", req.Email)
		return nil, errors.New("User already exists")
	}
	var user entity.User
	utils.Copy(&user, &req)
	err = s.repo.Create(ctx, &user)
	if err != nil {
		logger.Errorf("Register.Create fail, email: %s, error: %s", req.Email, err)
		return nil, err
	}
	return &user, nil
}

// func (s *UserService) GetMe(ctx context.Context){
// 	return nil
// }

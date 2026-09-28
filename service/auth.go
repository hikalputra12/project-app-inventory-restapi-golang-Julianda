package service

import (
	"app-inventory/dto"
	"app-inventory/repository"
	"app-inventory/utils"
	"context"
	"errors"

	"go.uber.org/zap"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
}

type authService struct {
	Repo      repository.UserRepoInterface
	Config    utils.JWTConfig
	logger    *zap.Logger
}

func NewAuthService(repo repository.UserRepoInterface, jwtConfig utils.JWTConfig, log *zap.Logger) AuthServiceInterface {
	return &authService{
		Repo:   repo,
		Config: jwtConfig,
		logger: log,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	user, err := s.Repo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("Login attempt failed: user email not found",
			zap.String("email", email))
		return nil, errors.New("email atau password salah")
	}

	// Never log plain text password!
	if !utils.CompareHashAndPassword(user.Password, password) {
		s.logger.Warn("Login attempt failed: incorrect password",
			zap.String("email", email))
		return nil, errors.New("email atau password salah")
	}

	// Generate JWT Access Token
	token, err := utils.GenerateJWT(user, s.Config.Secret, s.Config.ExpiryMinutes)
	if err != nil {
		s.logger.Error("Failed to generate JWT token", zap.Error(err))
		return nil, errors.New("failed to generate access token")
	}

	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		User: dto.UserProfile{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Role:   user.Role,
			RoleID: user.Role_id,
		},
	}, nil
}

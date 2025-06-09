package service

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/dto"
	"keeper/internal/entity"
	"keeper/internal/logger"
	"keeper/internal/repository"
	"time"

	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, requestDto *dto.LoginUser) (entity.AccessToken, error)
	Register(ctx context.Context, requestDto *dto.RegisterUser) (entity.AccessToken, error)
}

type authService struct {
	UserRepo   repository.UserRepositoryInterface
	AccessRepo repository.AccessTokenRepository
	JwtService JwtService
	db         *pgxpool.Pool
	l          *logger.ZapLogger
	cfg        config.SecurityConfig
}

func NewAuthService(
	db *pgxpool.Pool,
	userRepo repository.UserRepositoryInterface,
	accessRepo repository.AccessTokenRepository,
	jwtService JwtService,
	cfg config.SecurityConfig,
	l *logger.ZapLogger,
) AuthService {
	return &authService{
		db:         db,
		UserRepo:   userRepo,
		AccessRepo: accessRepo,
		JwtService: jwtService,
		cfg:        cfg,
		l:          l,
	}
}

const errCreateToken = "failed to create token: %w"

func (a *authService) Login(ctx context.Context, requestDto *dto.LoginUser) (entity.AccessToken, error) {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return entity.AccessToken{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
	}(tx, ctx)

	user, err := a.UserRepo.GetByLogin(ctx, tx, requestDto.Login)
	if err != nil {
		return entity.AccessToken{}, fmt.Errorf("failed to GetByLogin: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestDto.Password))
	if err != nil {
		return entity.AccessToken{}, fmt.Errorf("invalid credentials: %w", err)
	}

	token, err := a.AccessRepo.FindValidByUserID(ctx, tx, user.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return entity.AccessToken{}, fmt.Errorf("failed to check token: %w", err)
		}

		tokenString, err := a.JwtService.CreateJwt(user.ID)
		if err != nil {
			return entity.AccessToken{}, fmt.Errorf(errCreateToken, err)
		}
		newToken := &entity.AccessToken{
			UserID:    user.ID,
			Token:     tokenString,
			ExpiresAt: time.Now().Add(a.cfg.TokenTTL),
		}

		if err := a.AccessRepo.Create(ctx, tx, newToken); err != nil {
			return entity.AccessToken{}, fmt.Errorf(errCreateToken, err)
		}

		token = *newToken
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.AccessToken{}, fmt.Errorf("failed to commit tx: %w", err)
	}

	return token, nil
}

func (a *authService) Register(ctx context.Context, requestDto *dto.RegisterUser) (entity.AccessToken, error) {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return entity.AccessToken{}, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
	}(tx, ctx)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(requestDto.Password), bcrypt.DefaultCost)
	user := entity.User{Login: requestDto.Login, Password: string(hashedPassword)}

	newUser, err := a.UserRepo.Register(ctx, tx, user)
	if err != nil {
		return entity.AccessToken{}, fmt.Errorf("failed to register user: %w", err)
	}

	tokenString, err := a.JwtService.CreateJwt(newUser.ID)
	if err != nil {
		return entity.AccessToken{}, fmt.Errorf(errCreateToken, err)
	}

	token := &entity.AccessToken{
		UserID:    newUser.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(a.cfg.TokenTTL),
	}
	if err := a.AccessRepo.Create(ctx, tx, token); err != nil {
		return *token, fmt.Errorf(errCreateToken, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return *token, fmt.Errorf("failed to commit tx: %w", err)
	}

	return *token, nil
}

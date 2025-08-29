package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/entity"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/repository"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/utils"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
}

type authService struct {
	authRepository repository.IAUthRepository
}

func (as *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.Password != req.PasswordConfirmation {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password is not match"),
		}, nil
	}
	// cek email ke db
	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	// apabila email tedaftar, return error
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("User already exist"),
		}, nil
	}
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	// insert ke db
	newUser := entity.User{
		Id:        uuid.NewString(),
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		RoleCode:  entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: &req.FullName,
	}

	err = as.authRepository.InsertUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("User is registered"),
	}, nil
}

// Login implements IAuthService.
func (as *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	/// check apakah email ada
	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return  nil, err
	}

	if user == nil {
		return  &auth.LoginResponse{
			Base: utils.BadRequestResponse("User is not registered"),
		}, nil
	}

	// check apakah passwword sama
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
		}
		return nil, err
	}
	// generate jwt
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: user.Id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour *24)),
			IssuedAt: jwt.NewNumericDate(now),
		},
		Email: user.Email,
		FullName: user.FullName,
		Role: user.RoleCode,
	})

	secretKey := os.Getenv("JWT_SECRET")
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	// kirim response
	return &auth.LoginResponse{
		Base: utils.SuccessResponse("Login success"),
		AccessToken: accessToken,
	}, nil
}

// factory
func NewAuthService(authRepository repository.IAUthRepository) IAuthService {
	return &authService{
		authRepository: authRepository,
	}
}

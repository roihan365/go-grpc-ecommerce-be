package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	gocache "github.com/patrickmn/go-cache"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/entity"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/repository"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/utils"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	jwtEntity "github.com/roihan365/go-grpc-ecommerce-be/internal/entity/jwt"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
	GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error)
}

type authService struct {
	authRepository repository.IAUthRepository
	cacheService   *gocache.Cache
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
		return nil, err
	}

	if user == nil {
		return &auth.LoginResponse{
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtEntity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.RoleCode,
	})

	secretKey := os.Getenv("JWT_SECRET")
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	// kirim response
	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login success"),
		AccessToken: accessToken,
	}, nil
}

// Logout implements IAuthService.
func (as *authService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// get token dari metadata
	jwtToken, err := jwtEntity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// kembalikan token hingga menjadi entity jwt
	tokenClaim, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// masukkan token ke dalam memory db / cache
	as.cacheService.Set(jwtToken, "", time.Duration(tokenClaim.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)

	// kirim response
	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout success"),
	}, nil
}

func (as *authService) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	// get data token
	claims, err := jwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// ambil data dari db
	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return  &auth.GetProfileResponse{
			Base: utils.BadRequestResponse("User not found"),
		}, nil
	}
	// buat response

	// kirim response
	return &auth.GetProfileResponse{
		Base: utils.SuccessResponse("Get Profile success"),
		UserId: user.Id,
		FullName: user.FullName,
		Email: user.Email,
		RoleCode: user.RoleCode,
		MemberSince: timestamppb.New(user.CreatedAt),
	}, nil
}

// factory
func NewAuthService(authRepository repository.IAUthRepository, cacheService *gocache.Cache) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:  cacheService,
	}
}

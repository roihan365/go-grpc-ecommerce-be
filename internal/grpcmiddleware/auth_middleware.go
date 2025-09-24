package grpcmiddleware

import (
	"context"

	gocache "github.com/patrickmn/go-cache"
	jwtEntity "github.com/roihan365/go-grpc-ecommerce-be/internal/entity/jwt"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/utils"
	"google.golang.org/grpc"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if info.FullMethod == "/auth.AuthService/Login" || info.FullMethod == "/auth.AuthService/Register" {
		return handler(ctx, req)

	}

	// ambil token dari metadata
	jwtToken, err := jwtEntity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// cek token dari logout cache
	_, ok := am.cacheService.Get(jwtToken)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	// parse jwt hingga jadi entity
	claims, err := jwtEntity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}
	
	// sematkan entity ke context
	ctx = claims.SetToContext(ctx)
	res, err := handler(ctx, req)

	return res, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}
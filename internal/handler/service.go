package handler

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"github.com/roihan365/go-grpc-ecommerce-be/internal/utils"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serviceHandler struct {
	service.UnimplementedHelloWorldServer
}

func (sh *serviceHandler) HelloWorld(ctx context.Context, req *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error %v", err)
	}
	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s", req.Name),
		Base:    utils.SuccessResponse(),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}

package handler

import (
	"context"
	"fmt"

	"github.com/roihan365/go-grpc-ecommerce-be/pb/service"
)

type serviceHandler struct {
	service.UnimplementedHelloWorldServer
}

func (sh *serviceHandler) HelloWorld(ctx context.Context, req *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {
	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}

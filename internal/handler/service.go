package handler

import (
	"context"
	"fmt"

	"github.com/roihan365/go-grpc-ecommerce-be/internal/utils"
	"github.com/roihan365/go-grpc-ecommerce-be/pb/service"
)

type serviceHandler struct {
	service.UnimplementedHelloWorldServer
}

func (sh *serviceHandler) HelloWorld(ctx context.Context, req *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {
	validationErrors, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validationErrors != nil {
		return &service.HelloWorldResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	return &service.HelloWorldResponse{
		Message: fmt.Sprintf("Hello %s", req.Name),
		Base:    utils.SuccessResponse("Request success"),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}

package utils

import (
	"github.com/roihan365/go-grpc-ecommerce-be/pb/common"
)

func SuccessResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		Code:    200,
		Message: message,
	}
}

func BadRequestResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		Code:    400,
		Message: message,
		IsError: true,
	}
}

func ValidationErrorResponse(validationErrors []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		Code:            400,
		Message:         "Validation Error",
		IsError:         true,
		ValidationError: validationErrors,
	}
}

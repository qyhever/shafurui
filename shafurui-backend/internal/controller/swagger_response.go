package controller

import "shafurui/internal/model"

type SwaggerErrorResponse struct {
	Code    MyCode      `json:"code" example:"1001"`
	Message string      `json:"message" example:"请求参数错误"`
	Data    interface{} `json:"data"`
}

type SwaggerHelloInfoResponse struct {
	Code    MyCode                     `json:"code" example:"1000"`
	Message string                     `json:"message" example:"success"`
	Data    model.GetHelloInfoResponse `json:"data"`
}

type SwaggerAuthLoginResponse struct {
	Code    MyCode                  `json:"code" example:"1000"`
	Message string                  `json:"message" example:"success"`
	Data    model.AuthLoginResponse `json:"data"`
}

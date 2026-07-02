package repository

import (
	"shafurui/internal/model"
)

type AppRepository interface {
	GetHelloInfo(param *model.GetHelloInfoRequest) (*model.GetHelloInfoResponse, error)
}

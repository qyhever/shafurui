package persistence

import (
	"shafurui/internal/model"
	"shafurui/internal/repository"
)

type AppRepositoryImpl struct{}

func NewAppRepository() repository.AppRepository {
	return &AppRepositoryImpl{}
}

func (r *AppRepositoryImpl) GetHelloInfo(req *model.GetHelloInfoRequest) (*model.GetHelloInfoResponse, error) {
	res := &model.GetHelloInfoResponse{
		Name: req.Name,
	}
	return res, nil
}

package services

import (
	"gitlab.com/gjerry134679/bank/internal/dao"
	dbengines "gitlab.com/gjerry134679/bank/internal/models/dbEngines"
)

var DefaultService *Service

type Service struct {
	doa *dao.Dao
}

func NewService(DBsource dbengines.DBSource) (*Service, error) {
	doa, err := dao.New(DBsource)
	if err != nil {
		return nil, err
	}
	return &Service{doa: doa}, nil
}

func AsignDefaultService(srv *Service) {
	DefaultService = srv
}

package dao

import (

	// "github.com/jackc/pgx-zerolog"

	_ "github.com/jackc/pgx/v4"
	"gitlab.com/gjerry134679/bank/internal/models"
	dbengines "gitlab.com/gjerry134679/bank/internal/models/dbEngines"
)

var DefaultDao Dao

type Dao struct {
	engine *models.Queries
}

func New(dbSource dbengines.DBSource) (*Dao, error) {
	engine, err := dbengines.NewEngine(dbSource)
	if err != nil {
		return nil, err
	}
	return &Dao{engine: models.New(engine)}, nil
}

func AsignDefaultDao(dao Dao) {
	DefaultDao = dao
}

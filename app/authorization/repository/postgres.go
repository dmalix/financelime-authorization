package repository

import (
	"database/sql"
	"fmt"
	"github.com/dmalix/financelime-authorization/app/authorization/model"
	"go.uber.org/zap"
)

func NewPostgresDB(logger *zap.Logger, c model.ConfigPostgresDB) (*sql.DB, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d sslmode=%s dbname=%s user=%s password=%s",
			c.Host, c.Port, c.SSLMode, c.DBName, c.User, c.Password))
	if err != nil {
		logger.DPanic("failed to open PostgresDB connection", zap.Error(err))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.DPanic("failed to ping for PostgresDB", zap.Error(err))
		return nil, err
	}

	return db, nil
}

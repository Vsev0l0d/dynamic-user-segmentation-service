package db

import (
	"dynamic-user-segmentation-service/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

func Connect(dbConfig config.DB) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, os.Getenv(dbConfig.EnvUser), os.Getenv(dbConfig.EnvPassword),
		dbConfig.DbName, dbConfig.SslMode)
	db, err := sqlx.Connect(dbConfig.DriverName, dataSourceName)
	return db, err
}

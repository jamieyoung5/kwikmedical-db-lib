package config

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"net/url"
)

const driver = "postgresql"

func CreateSqlClient(logger *zap.Logger, dbConfig *Config) (*sql.DB, error) {
	connectionStr := createConnectionString(dbConfig)
	logger.Debug(connectionStr)

	client, err := sql.Open("postgres", connectionStr)

	if err != nil {
		logger.Error("Error opening database", zap.Error(err))
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		logger.Error("Failed to Ping() using client", zap.Error(err))
		return nil, err
	}

	logger.Info("Successfully pinged database")

	return client, nil
}

func createConnectionString(dbConfig *Config) string {
	//postgresql://username:password@host/database?sslmode=require
	return fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=require",
		driver,
		dbConfig.UserName,
		url.QueryEscape(dbConfig.Password),
		dbConfig.Host,
		dbConfig.DatabaseName)
}

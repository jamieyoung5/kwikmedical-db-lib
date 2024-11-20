package config

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"net/url"
)

func CreateSqlClient(logger *zap.Logger, dbConfig *Config) (*sql.DB, error) {
	connectionStr := createConnectionString(dbConfig)
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

	logger.Info("Successfully pinged database.")

	return client, nil
}

func createConnectionString(dbConfig *Config) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=require",
		dbConfig.UserName,
		url.QueryEscape(dbConfig.Password),
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DatabaseName)
}

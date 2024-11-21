package client

import (
	"database/sql"
	"github.com/jamieyoung5/quickmedical-db-lib/pkg/config"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SqlDb interface {
	gorm.ConnPool
	Ping() error
	Close() error
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type KwikMedicalDBClient struct {
	logger      *zap.Logger
	gormDb      *gorm.DB
	sqlDb       SqlDb
	isConnected bool
}

func NewKwikMedicalDBClient(logger *zap.Logger, gormDb *gorm.DB) (*KwikMedicalDBClient, error) {

	sqlDb, err := gormDb.DB()
	if err != nil {
		logger.Error("Error connecting to SQL database", zap.Error(err))
		return nil, err
	}

	return &KwikMedicalDBClient{
		logger: logger,
		gormDb: gormDb,
		sqlDb:  sqlDb,
	}, nil
}

func NewClient(logger *zap.Logger, dbConfig *config.Config) (*KwikMedicalDBClient, error) {
	sqlDb, err := config.CreateSqlClient(logger, dbConfig)
	if err != nil {
		logger.Error("Error creating SQL client", zap.Error(err))
		return nil, err
	}

	gormDb, err := gorm.Open(
		postgres.New(
			postgres.Config{
				Conn: sqlDb,
			}),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	return NewKwikMedicalDBClient(logger, gormDb)
}

func (db *KwikMedicalDBClient) IsConnected() bool {
	return db.isConnected
}

func (db *KwikMedicalDBClient) Ping() error {
	err := db.sqlDb.Ping()
	if err != nil {
		db.logger.Error("Failed to Ping", zap.Error(err))
		return err
	}

	db.logger.Debug("Successfully pinged database")
	db.isConnected = true

	return nil
}

func (db *KwikMedicalDBClient) Close() error {
	err := db.sqlDb.Close()
	if err != nil {
		db.logger.Error("Failed to Close", zap.Error(err))
		return err
	}

	db.logger.Debug("Successfully closed database")
	return nil
}

func (db *KwikMedicalDBClient) Exec(query string, args ...any) (sql.Result, error) {
	return db.sqlDb.Exec(query, args...)
}

func (db *KwikMedicalDBClient) Query(query string, args ...any) (*sql.Rows, error) {
	return db.sqlDb.Query(query, args...)
}

func (db *KwikMedicalDBClient) QueryRow(query string, args ...any) *sql.Row {
	return db.sqlDb.QueryRow(query, args...)
}

func (db *KwikMedicalDBClient) DbTransaction(fn func(tx *gorm.DB) error) error {
	tx := db.gormDb.Begin()
	if tx.Error != nil {
		db.logger.Error("Error starting transaction", zap.Error(tx.Error))
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := fn(tx)
	if err != nil {
		db.logger.Error("Error executing transaction operation", zap.Error(err))
		return err
	}

	tx.Commit()

	return nil
}

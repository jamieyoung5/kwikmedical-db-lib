package main

import (
	dbClient "github.com/jamieyoung5/quickmedical-db-lib/pkg/client"
	dbConfig "github.com/jamieyoung5/quickmedical-db-lib/pkg/config"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Error("Error initializing logger", zap.Error(err))
		os.Exit(1)
	}

	config := dbConfig.NewConfig()

	client, err := dbClient.NewClient(logger, config)
	if err != nil {
		logger.Error("Error initializing client", zap.Error(err))
		os.Exit(1)
	}

	callInfo := dbClient.EmergencyCallPatientInfo{
		Address:   "123 Main St, Anytown",
		LastName:  "Doe",
		FirstName: "John",
	}
	id, err := client.FindClosestPatientID(callInfo)
	if err != nil {
		logger.Error("Error finding patient ID", zap.Error(err))
		os.Exit(1)
	}

	patient, err := client.GetPatientByID(id)
	if err != nil {
		logger.Error("Error getting patient", zap.Error(err))
		os.Exit(1)
	}
	logger.Debug("got patient successfully", zap.Any("patient", patient))

	patientData, err := client.GetHistoricalPatientDataByID(id)
	if err != nil {
		logger.Error("Error getting patient data", zap.Error(err))
		os.Exit(1)
	}
	logger.Debug("got historical data successfully", zap.Any("historical data", patientData))

	err = client.Close()
	if err != nil {
		logger.Error("Error closing client connection", zap.Error(err))
		os.Exit(1)
	}
}

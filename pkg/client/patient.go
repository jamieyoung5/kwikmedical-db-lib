package client

import (
	"errors"
	"github.com/jamieyoung5/kwikmedical-db-lib/pkg/schema"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmergencyCallPatientInfo struct {
	FirstName string
	LastName  string
	Address   string
}

type HistoricalPatientData struct {
	Patient       *schema.Patient
	MedicalRecord *schema.MedicalRecord
	Callouts      []schema.CallOutDetails
}

func (db *KwikMedicalDBClient) GetHistoricalPatientDataByID(id uint) (HistoricalPatientData, error) {
	patient, err := db.GetPatientByID(id)
	if err != nil {
		db.logger.Error("Unable to get patient", zap.Int("id", int(id)), zap.Error(err))
		return HistoricalPatientData{}, err
	}

	medicalRecord, callouts, err := db.GetMedicalRecordsByPatientID(id)
	if err != nil {
		db.logger.Error("Unable to get medical records", zap.Int("id", int(id)), zap.Error(err))
		return HistoricalPatientData{Patient: patient}, err
	}

	return HistoricalPatientData{
		Patient:       patient,
		MedicalRecord: medicalRecord,
		Callouts:      callouts,
	}, nil
}

func (db *KwikMedicalDBClient) GetPatientByID(id uint) (*schema.Patient, error) {
	var patient schema.Patient

	err := db.DbTransaction(func(tx *gorm.DB) error {
		if err := tx.First(&patient, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("patient not found")
			}
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &patient, nil
}

func (db *KwikMedicalDBClient) FindClosestPatientID(callInfo EmergencyCallPatientInfo) (uint, error) {
	var patient schema.Patient

	// construct search clause
	var searchClause []clause.Expression
	if callInfo.FirstName != "" {
		searchClause = append(searchClause, clause.Eq{
			Column: clause.Column{Name: "first_name"},
			Value:  callInfo.FirstName,
		})
	}
	if callInfo.LastName != "" {
		searchClause = append(searchClause, clause.Eq{
			Column: clause.Column{Name: "last_name"},
			Value:  callInfo.LastName,
		})
	}
	if callInfo.Address != "" {
		searchClause = append(searchClause, clause.Eq{
			Column: clause.Column{Name: "address"},
			Value:  callInfo.Address,
		})
	}

	// tries combinations of name and address to find the best patient match
	err := db.DbTransaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(searchClause...).
			First(&patient).Error; err == nil {
			return nil
		}

		return errors.New("patient not found")
	})

	if err != nil {
		return 0, err
	}

	return patient.PatientID, nil
}

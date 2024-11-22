package client

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EmergencyCallPatientInfo struct {
	FirstName string
	LastName  string
	Address   string
}

type HistoricalPatientData struct {
	Patient       *Patient
	MedicalRecord *MedicalRecord
	Callouts      []CallOutDetails
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

func (db *KwikMedicalDBClient) GetPatientByID(id uint) (*Patient, error) {
	var patient Patient

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
	var patient Patient

	// tries combinations of name and address to find the best patient match
	err := db.DbTransaction(func(tx *gorm.DB) error {
		if err := tx.Where(
			"first_name = ? AND last_name = ? AND address = ?",
			callInfo.FirstName,
			callInfo.LastName,
			callInfo.Address).
			First(&patient).Error; err == nil {
			return nil
		}

		if err := tx.Where(
			"first_name = ? AND last_name = ?",
			callInfo.FirstName,
			callInfo.LastName).
			First(&patient).Error; err == nil {
			return nil
		}

		if err := tx.Where(
			"last_name = ? AND address = ?",
			callInfo.LastName,
			callInfo.Address).
			First(&patient).Error; err == nil {
			return nil
		}

		if err := tx.Where(
			"first_name = ? AND address = ?",
			callInfo.FirstName,
			callInfo.Address).
			First(&patient).Error; err == nil {
			return nil
		}

		if err := tx.Where(
			"last_name = ?",
			callInfo.LastName).
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

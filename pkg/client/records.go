package client

import (
	"errors"
	"fmt"
	"github.com/jamieyoung5/kwikmedical-db-lib/pkg/schema"
	"github.com/jamieyoung5/kwikmedical-eventstream/pb"
	"gorm.io/gorm"
)

func (db *KwikMedicalDBClient) InsertNewCallout(callout *pb.CallOutDetail) error {
	calloutDetails := schema.CalloutDetailPbToGorm(callout)

	err := db.gormDb.Create(&calloutDetails).Error
	if err != nil {
		return err
	}

	patientId, err := db.GetPatientByEmergencyCall(calloutDetails.CallID)
	err = db.gormDb.Exec(
		`UPDATE medical_records SET callout_ids = array_append(callout_ids, ?) WHERE patient_id = ?`,
		calloutDetails.DetailID,
		patientId).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *KwikMedicalDBClient) GetMedicalRecordsByEmergencyCall(id uint) (*schema.MedicalRecord, []schema.CallOutDetails, error) {
	patientId, err := db.GetPatientByEmergencyCall(id)
	if err != nil {
		return nil, nil, err
	}

	return db.GetMedicalRecordsByPatientID(patientId)
}

func (db *KwikMedicalDBClient) GetMedicalRecordsByPatientID(id uint) (*schema.MedicalRecord, []schema.CallOutDetails, error) {
	var (
		medicalRecord  schema.MedicalRecord
		callOutDetails []schema.CallOutDetails
	)

	err := db.DbTransaction(func(tx *gorm.DB) error {
		if err := tx.Where("patient_id = ?", id).
			Order("last_updated DESC").
			First(&medicalRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("no medical records found for patient_id %d", id)
			}
			return err
		}

		if medicalRecord.RecordID == 0 {
			return fmt.Errorf("no medical records found for patient ID %d", id)
		}

		if err := tx.Where("detail_id IN ?", []int64(medicalRecord.CalloutIDs)).Find(&callOutDetails).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("no callout details found for the given IDs")
			}
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &medicalRecord, callOutDetails, nil
}

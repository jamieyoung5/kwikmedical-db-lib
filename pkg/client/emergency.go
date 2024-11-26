package client

import (
	"github.com/jamieyoung5/kwikmedical-db-lib/pkg/schema"
	"github.com/jamieyoung5/kwikmedical-eventstream/pb"
	"gorm.io/gorm"
)

func (db *KwikMedicalDBClient) AssignAmbulance(requestId int) (*int32, error) {
	var ambulanceID *int32
	err := db.DbTransaction(func(tx *gorm.DB) error {
		err := tx.Table("ambulances").
			Select("ambulances.ambulance_id").
			Joins("INNER JOIN ambulance_requests ON ambulances.regional_hospital_id = ambulance_requests.hospital_id").
			Where("ambulance_requests.request_id = ?", requestId).
			Where("ambulances.status = ?", "AVAILABLE").
			Limit(1).
			Scan(&ambulanceID).Error

		if err != nil {
			return err
		}

		if ambulanceID == nil {
			db.logger.Debug("no ambulances to assign")
			return nil
		}

		err = tx.Table("ambulance_requests").
			Where("request_id = ?", requestId).
			Update("ambulance_id", *ambulanceID).Error
		if err != nil {
			return err
		}

		err = tx.Table("ambulance_requests").
			Where("request_id = ?", requestId).
			Update("status", "ACCEPTED").Error
		if err != nil {
			return err
		}

		err = tx.Table("ambulances").
			Where("ambulance_id = ?", *ambulanceID).
			Update("status", "ON_CALL").Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return ambulanceID, nil
}

func (db *KwikMedicalDBClient) CreateNewAmbulanceRequest(request *pb.AmbulanceRequest) (int32, error) {
	ambulanceRequest := schema.AmbulanceRequestPbToGorm(request)

	if err := db.gormDb.Create(&ambulanceRequest).Error; err != nil {
		return 0, err
	}

	return int32(ambulanceRequest.RequestID), nil
}

func (db *KwikMedicalDBClient) InsertNewEmergencyCall(call *pb.EmergencyCall) (int32, error) {
	emergencyCall := schema.EmergencyCallPbToGorm(call)

	if err := db.gormDb.Create(&emergencyCall).Error; err != nil {
		return 0, err
	}

	return int32(emergencyCall.CallID), nil
}

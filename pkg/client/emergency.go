package client

import "github.com/jamieyoung5/kwikmedical-eventstream/pb"

func (db *KwikMedicalDBClient) InsertNewEmergencyCall(call *pb.EmergencyCall) error {
	emergencyCall := EmergencyCallPbToGorm(call)

	if err := db.gormDb.Create(emergencyCall).Error; err != nil {
		return err
	}

	return nil
}

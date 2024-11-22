package client

import "github.com/jamieyoung5/kwikmedical-eventstream/pb"

func (db *KwikMedicalDBClient) InsertNewEmergencyCall(call *pb.EmergencyCall) (int32, error) {
	emergencyCall := EmergencyCallPbToGorm(call)

	if err := db.gormDb.Create(emergencyCall).Error; err != nil {
		return 0, err
	}

	return emergencyCall.CallId, nil
}

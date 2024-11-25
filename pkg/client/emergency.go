package client

import (
	"github.com/jamieyoung5/kwikmedical-db-lib/pkg/schema"
	"github.com/jamieyoung5/kwikmedical-eventstream/pb"
)

func (db *KwikMedicalDBClient) CreateNewAmbulanceRequest(request *pb.AmbulanceRequest) (int32, error) {
	ambulanceRequest := schema.AmbulanceRequestPbToGorm(request)

	if request.HospitalId == 0 {
		hospital, err := db.GetNearestHospital(request.Location)
		if err != nil {
			return 0, err
		}
		ambulanceRequest.HospitalID = &hospital.HospitalID
	}

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

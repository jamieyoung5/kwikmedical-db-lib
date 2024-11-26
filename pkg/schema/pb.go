package schema

import pbSchema "github.com/jamieyoung5/kwikmedical-eventstream/pb"

func EmergencyCallPbToGorm(call *pbSchema.EmergencyCall) EmergencyCall {
	patientId := uint(call.PatientId)

	return EmergencyCall{
		CallID:           uint(call.CallId),
		PatientID:        &patientId,
		NHSNumber:        call.NhsNumber,
		CallerName:       call.CallerName,
		CallerPhone:      call.CallerPhone,
		CallTime:         call.CallTime.AsTime(),
		MedicalCondition: call.MedicalCondition,
		Location:         LocationFromPb(call.Location),
		Severity:         InjurySeverity(pbSchema.InjurySeverity_name[int32(call.Severity)]),
		Status:           EmergencyCallStatus(pbSchema.EmergencyCallStatus_name[int32(call.Status)]),
	}
}

func AmbulanceRequestPbToGorm(request *pbSchema.AmbulanceRequest) AmbulanceRequest {
	hospitalId := uint(request.HospitalId)

	return AmbulanceRequest{
		RequestID:       uint(request.RequestId),
		HospitalID:      &hospitalId,
		EmergencyCallID: uint(request.EmergencyCallId),
		Severity:        InjurySeverity(pbSchema.InjurySeverity_name[int32(request.Severity)]),
		Location:        LocationFromPb(request.Location),
		Status:          RequestStatus(pbSchema.RequestStatus_name[int32(request.Status)]),
		CreatedAt:       request.CreatedAt.AsTime(),
		UpdatedAt:       request.UpdatedAt.AsTime(),
	}
}

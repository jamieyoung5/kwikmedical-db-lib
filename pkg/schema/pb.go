package schema

import pbSchema "github.com/jamieyoung5/kwikmedical-eventstream/pb"

func EmergencyCallPbToGorm(call *pbSchema.EmergencyCall) EmergencyCall {
	patientId := uint(call.PatientId)
	assignedAmbulanceId := uint(call.AssignedAmbulanceId)
	assignedHospitalId := uint(call.AssignedHospitalId)

	return EmergencyCall{
		CallID:              uint(call.CallId),
		PatientID:           &patientId,
		NHSNumber:           call.NhsNumber,
		CallerName:          call.CallerName,
		CallerPhone:         call.CallerPhone,
		CallTime:            call.CallTime.AsTime(),
		MedicalCondition:    call.MedicalCondition,
		Location:            PointFromPb(call.Location),
		Severity:            InjurySeverity(pbSchema.InjurySeverity_name[int32(call.Severity)]),
		Status:              EmergencyCallStatus(pbSchema.EmergencyCallStatus_name[int32(call.Status)]),
		AssignedAmbulanceID: &assignedAmbulanceId,
		AssignedHospitalID:  &assignedHospitalId,
	}
}

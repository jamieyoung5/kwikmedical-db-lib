package schema

type EmergencyCallStatus string
type AmbulanceStatus string
type InjurySeverity string
type StaffRole string
type RequestStatus string

const (
	UnknownEmergency EmergencyCallStatus = "UNKNOWN_EMERGENCY_CALL_STATUS"
	Pending          EmergencyCallStatus = "AMBULANCE_PENDING"
	Dispatched       EmergencyCallStatus = "AMBULANCE_DISPATCHED"
	Completed        EmergencyCallStatus = "AMBULANCE_COMPLETED"

	UnknownAmbulance AmbulanceStatus = "UNKNOWN_AMBULANCE_STATUS"
	Available        AmbulanceStatus = "AVAILABLE"
	OnCall           AmbulanceStatus = "ON_CALL"
	Maintenance      AmbulanceStatus = "MAINTENANCE"

	UnknownSeverity InjurySeverity = "UNKNOWN_INJURY_SEVERITY"
	Low             InjurySeverity = "LOW"
	Moderate        InjurySeverity = "MODERATE"
	High            InjurySeverity = "HIGH"
	Critical        InjurySeverity = "CRITICAL"

	UnknownRole   StaffRole = "UNKNOWN_STAFF_ROLE"
	Paramedic     StaffRole = "PARAMEDIC"
	Driver        StaffRole = "DRIVER"
	Operator      StaffRole = "OPERATOR"
	HospitalStaff StaffRole = "HOSPITAL_STAFF"
	Other         StaffRole = "OTHER"

	UnknownReq   RequestStatus = "UNKNOWN_REQUEST_STATUS"
	ReqPending   RequestStatus = "PENDING"
	ReqAccepted  RequestStatus = "ACCEPTED"
	ReqRejected  RequestStatus = "REJECTED"
	ReqCompleted RequestStatus = "COMPLETED"
)

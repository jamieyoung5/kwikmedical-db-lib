package schema

import (
	pbSchema "github.com/jamieyoung5/kwikmedical-eventstream/pb"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Patient struct {
	PatientID   uint      `gorm:"primaryKey" json:"patient_id"`
	NHSNumber   string    `gorm:"type:varchar(15);unique;not null" json:"nhs_number"`
	FirstName   string    `gorm:"type:varchar(50);not null" json:"first_name"`
	LastName    string    `gorm:"type:varchar(50);not null" json:"last_name"`
	DateOfBirth string    `gorm:"type:date" json:"date_of_birth"`
	Address     string    `gorm:"type:text" json:"address"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phone_number"`
	Email       string    `gorm:"type:varchar(100)" json:"email"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (p *Patient) ToPb() *pbSchema.Patient {
	var createdAtPb *timestamppb.Timestamp
	if !p.CreatedAt.IsZero() {
		createdAtPb = timestamppb.New(p.CreatedAt)
	}

	return &pbSchema.Patient{
		PatientId:   int32(p.PatientID),
		NhsNumber:   p.NHSNumber,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		DateOfBirth: p.DateOfBirth,
		Address:     p.Address,
		PhoneNumber: p.PhoneNumber,
		Email:       p.Email,
		CreatedAt:   createdAtPb,
	}
}

type MedicalRecord struct {
	RecordID    uint           `gorm:"primaryKey" json:"record_id"`
	PatientID   uint           `gorm:"not null;constraint:OnDelete:CASCADE" json:"patient_id"`
	CalloutIDs  pq.Int64Array  `gorm:"type:int[]" json:"callout_ids"` // For array, PostgreSQL specific
	Conditions  pq.StringArray `gorm:"type:text[]" json:"conditions"`
	Medications pq.StringArray `gorm:"type:text[]" json:"medications"`
	Allergies   pq.StringArray `gorm:"type:text[]" json:"allergies"`
	Notes       pq.StringArray `gorm:"type:text[]" json:"notes"`
	LastUpdated time.Time      `gorm:"autoUpdateTime" json:"last_updated"`
}

func (mr *MedicalRecord) ToPb(callouts []CallOutDetails) *pbSchema.MedicalRecord {
	var calloutDetails []*pbSchema.CallOutDetail
	for _, callout := range callouts {
		calloutDetails = append(calloutDetails, callout.ToPb())
	}

	return &pbSchema.MedicalRecord{
		RecordId:    int32(mr.RecordID),
		Callouts:    calloutDetails,
		Conditions:  mr.Conditions,
		Medications: mr.Medications,
		Allergies:   mr.Allergies,
		Notes:       mr.Notes,
		LastUpdated: timestamppb.New(mr.LastUpdated),
	}
}

type CallOutDetails struct {
	DetailID    uint      `gorm:"primaryKey" json:"detail_id"`
	CallID      uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"call_id"`
	AmbulanceID uint      `json:"ambulance_id"`
	ActionTaken string    `gorm:"type:text" json:"action_taken"`
	TimeSpent   string    `gorm:"type:interval" json:"time_spent"`
	Notes       string    `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (cd *CallOutDetails) ToPb() *pbSchema.CallOutDetail {
	var timeSpent *durationpb.Duration
	duration, err := time.ParseDuration(cd.TimeSpent)
	if err == nil {
		timeSpent = durationpb.New(duration)
	}

	return &pbSchema.CallOutDetail{
		DetailId:    int32(cd.DetailID),
		CallId:      int32(cd.CallID),
		AmbulanceId: int32(cd.AmbulanceID),
		ActionTaken: cd.ActionTaken,
		TimeSpent:   timeSpent,
		Notes:       cd.Notes,
		CreatedAt:   timestamppb.New(cd.CreatedAt),
	}
}

type EmergencyCall struct {
	CallID              uint                `gorm:"primaryKey;autoIncrement" json:"call_id"`
	PatientID           *uint               `gorm:"constraint:OnDelete:SET NULL" json:"patient_id"`
	NHSNumber           string              `gorm:"type:varchar(15)" json:"nhs_number"`
	CallerName          string              `gorm:"type:varchar(100)" json:"caller_name"`
	CallerPhone         string              `gorm:"type:varchar(20)" json:"caller_phone"`
	CallTime            time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"call_time"`
	MedicalCondition    string              `gorm:"type:text" json:"medical_condition"`
	Location            Point               `gorm:"type:text" json:"location"`
	Severity            InjurySeverity      `gorm:"type:injury_severity;default:'Low'" json:"severity"`
	Status              EmergencyCallStatus `gorm:"type:emergency_call_status;default:'Pending'" json:"status"`
	AssignedAmbulanceID *uint               `gorm:"constraint:OnDelete:SET NULL" json:"assigned_ambulance_id"`
	AssignedHospitalID  *uint               `gorm:"constraint:OnDelete:SET NULL" json:"assigned_hospital_id"`
}

type Ambulance struct {
	AmbulanceID        uint            `gorm:"primaryKey" json:"ambulance_id"`
	AmbulanceNumber    string          `gorm:"type:varchar(20);unique;not null" json:"ambulance_number"`
	CurrentLocation    Point           `gorm:"type:point" json:"current_location"` // PostGIS POINT type
	Status             AmbulanceStatus `gorm:"type:ambulance_status;default:'Available'" json:"status"`
	RegionalHospitalID *uint           `gorm:"constraint:OnDelete:SET NULL" json:"regional_hospital_id"`
}

type AmbulanceRequest struct {
	RequestID       uint           `gorm:"primaryKey;autoIncrement" json:"request_id"`
	AmbulanceID     uint           `gorm:"not null" json:"ambulance_id"`
	HospitalID      *uint          `json:"hospital_id"`
	EmergencyCallID uint           `gorm:"not null;constraint:OnDelete:CASCADE" json:"emergency_call_id"`
	Severity        InjurySeverity `gorm:"type:injury_severity" json:"severity"`
	Location        Point          `gorm:"type:point" json:"location"` // PostGIS POINT type
	Status          RequestStatus  `gorm:"type:request_status" json:"status"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type AmbulanceStaff struct {
	StaffID     uint      `gorm:"primaryKey" json:"staff_id"`
	FirstName   string    `gorm:"type:varchar(50);not null" json:"first_name"`
	LastName    string    `gorm:"type:varchar(50);not null" json:"last_name"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phone_number"`
	Email       string    `gorm:"type:varchar(100)" json:"email"`
	Role        StaffRole `gorm:"type:staff_role" json:"role"`
	AmbulanceID *uint     `json:"ambulance_id"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}

type RegionalHospital struct {
	HospitalID  uint      `gorm:"primaryKey" json:"hospital_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Address     string    `gorm:"type:text" json:"address"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phone_number"`
	Email       string    `gorm:"type:varchar(100)" json:"email"`
	Location    Point     `gorm:"type:point" json:"location"` // PostGIS POINT type
	Capacity    int       `json:"capacity"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (rh *RegionalHospital) ToPb() *pbSchema.RegionalHospital {
	return &pbSchema.RegionalHospital{
		HospitalId:  int32(rh.HospitalID),
		Name:        rh.Name,
		Address:     rh.Address,
		PhoneNumber: rh.PhoneNumber,
		Email:       rh.Email,
		Location: &pbSchema.Location{
			Latitude:  rh.Location.Latitude,
			Longitude: rh.Location.Longitude,
		},
		Capacity:  int32(rh.Capacity),
		CreatedAt: timestamppb.New(rh.CreatedAt),
	}
}

package client

import (
	pbSchema "github.com/jamieyoung5/kwikmedical-eventstream/pb"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Patient struct {
	PatientID   uint      `gorm:"primaryKey;column:patient_id"`
	NHSNumber   string    `gorm:"column:nhs_number"`
	FirstName   string    `gorm:"column:first_name"`
	LastName    string    `gorm:"column:last_name"`
	DateOfBirth string    `gorm:"column:date_of_birth"`
	Address     string    `gorm:"column:address"`
	PhoneNumber string    `gorm:"column:phone_number"`
	Email       string    `gorm:"column:email"`
	CreatedAt   time.Time `gorm:"column:created_at"`
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
	RecordID    int            `gorm:"primaryKey;autoIncrement;column:record_id"`
	PatientID   int            `gorm:"index;not null;column:patient_id"`
	Patient     *Patient       `gorm:"foreignKey:PatientID;constraint:OnDelete:CASCADE;references:PatientID"`
	CalloutIDs  pq.Int64Array  `gorm:"type:integer[];column:callout_ids"`
	Conditions  pq.StringArray `gorm:"type:text[];column:conditions"`
	Medications pq.StringArray `gorm:"type:text[];column:medications"`
	Allergies   pq.StringArray `gorm:"type:text[];column:allergies"`
	Notes       pq.StringArray `gorm:"type:text[];column:notes"`
	LastUpdated time.Time      `gorm:"autoUpdateTime;default:CURRENT_TIMESTAMP;column:last_updated"`
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
	DetailID    int       `gorm:"primaryKey;autoIncrement;column:detail_id"`
	CallID      int       `gorm:"column:call_id"`
	AmbulanceID int       `gorm:"column:ambulance_id"`
	ActionTaken string    `gorm:"column:action_taken"`
	TimeSpent   string    `gorm:"column:time_spent"`
	Notes       string    `gorm:"column:notes"`
	CreatedAt   time.Time `gorm:"column:created_at"`
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
	CallId              int32     `gorm:"primaryKey;column:call_id;autoIncrement"`
	PatientId           *int32    `gorm:"column:patient_id"`
	NhsNumber           string    `gorm:"column:nhs_number"`
	CallerName          string    `gorm:"column:caller_name"`
	CallerPhone         string    `gorm:"column:caller_phone"`
	CallTime            time.Time `gorm:"column:call_time"`
	MedicalCondition    string    `gorm:"column:medical_condition"`
	Location            string    `gorm:"column:location"`
	Status              string    `gorm:"column:status;default:Pending"`
	AssignedAmbulanceId *int32    `gorm:"column:assigned_ambulance_id"`
	AssignedHospitalId  *int32    `gorm:"column:assigned_hospital_id"`
}

func EmergencyCallPbToGorm(proto *pbSchema.EmergencyCall) *EmergencyCall {
	var callTime time.Time
	if proto.CallTime != nil {
		callTime = proto.CallTime.AsTime()
	}

	return &EmergencyCall{
		CallId:              proto.CallId,
		PatientId:           int32Pointer(proto.PatientId),
		NhsNumber:           proto.NhsNumber,
		CallerName:          proto.CallerName,
		CallerPhone:         proto.CallerPhone,
		CallTime:            callTime,
		MedicalCondition:    proto.MedicalCondition,
		Location:            proto.Location,
		Status:              proto.Status,
		AssignedAmbulanceId: int32Pointer(proto.AssignedAmbulanceId),
		AssignedHospitalId:  int32Pointer(proto.AssignedHospitalId),
	}
}

func int32Pointer(value int32) *int32 {
	return &value
}

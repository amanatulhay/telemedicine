package structs

import "time"

type Consultation struct {
	ID          int64     `json:"id"`
	MeetingLink string    `json:"meeting_link" binding:"required"`
	PaymentLink string    `json:"payment_link" binding:"required"`
	PatientID   int64     `json:"patient_id"`
	DoctorID    int64     `json:"doctor_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

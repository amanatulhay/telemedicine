package structs

import "time"

type Prescription struct {
	ID             int64     `json:"id"`
	Content        string    `json:"content" binding:"required"`
	PaymentLink    string    `json:"payment_link" binding:"required"`
	PatientID      int64     `json:"patient_id" binding:"required"`
	DoctorID       int64     `json:"doctor_id"`
	ConsultationID int64     `json:"consultation_id" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

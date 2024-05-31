package repository

import (
	"database/sql"
	"telemedicine/structs"
)

func GetAllConsultations(db *sql.DB) (err error, results []structs.Consultation) {
	sql := "SELECT * FROM consultation"

	rows, err := db.Query(sql)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var consultation = structs.Consultation{}

		err = rows.Scan(&consultation.ID, &consultation.MeetingLink, &consultation.PaymentLink, &consultation.PatientID, &consultation.DoctorID, &consultation.CreatedAt, &consultation.UpdatedAt)
		if err != nil {
			return
		}

		results = append(results, consultation)
	}

	return
}

func InsertConsultation(db *sql.DB, consultation structs.Consultation) (err error) {
	sql := "INSERT INTO consultation(id, meeting_link, payment_link, patient_id, doctor_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	errs := db.QueryRow(sql, consultation.ID, consultation.MeetingLink, consultation.PaymentLink, consultation.PatientID, consultation.DoctorID, consultation.CreatedAt, consultation.UpdatedAt)

	return errs.Err()
}

func UpdateConsultation(db *sql.DB, consultation structs.Consultation) (err error) {

	sql := "UPDATE consultation SET meeting_link = $1, payment_link = $2, patient_id = $3, doctor_id = $4, created_at = $5, updated_at = $6 WHERE id = $7"

	errs := db.QueryRow(sql, consultation.MeetingLink, consultation.PaymentLink, consultation.PatientID, consultation.DoctorID, consultation.CreatedAt, consultation.UpdatedAt, consultation.ID)

	return errs.Err()
}

func DeleteConsultation(db *sql.DB, consultation structs.Consultation) (err error) {
	sql := "DELETE FROM consultation WHERE id = $1"

	errs := db.QueryRow(sql, consultation.ID)

	return errs.Err()
}

func GetAllConsultationsByPatientID(db *sql.DB, patient structs.Patient) (err error, results []structs.Consultation) {
	sql := "SELECT * FROM consultation WHERE patient_id = $1"

	rows, err := db.Query(sql, patient.ID)
	if err != nil {
		return err, nil
	}

	defer rows.Close()

	for rows.Next() {
		var consultation = structs.Consultation{}

		err = rows.Scan(&consultation.ID, &consultation.MeetingLink, &consultation.PaymentLink, &consultation.PatientID, &consultation.DoctorID, &consultation.CreatedAt, &consultation.UpdatedAt)
		if err != nil {
			panic(err)
		}

		results = append(results, consultation)
	}

	return
}

func GetAllConsultationsByDoctorID(db *sql.DB, doctor structs.Doctor) (err error, results []structs.Consultation) {
	sql := "SELECT * FROM consultation WHERE doctor_id = $1"

	rows, err := db.Query(sql, doctor.ID)
	if err != nil {
		return err, nil
	}

	defer rows.Close()

	for rows.Next() {
		var consultation = structs.Consultation{}

		err = rows.Scan(&consultation.ID, &consultation.MeetingLink, &consultation.PaymentLink, &consultation.PatientID, &consultation.DoctorID, &consultation.CreatedAt, &consultation.UpdatedAt)
		if err != nil {
			panic(err)
		}

		results = append(results, consultation)
	}

	return
}

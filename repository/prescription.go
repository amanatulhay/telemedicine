package repository

import (
	"database/sql"
	"telemedicine/structs"
)

func GetAllPrescriptions(db *sql.DB) (err error, results []structs.Prescription) {
	sql := "SELECT * FROM prescription"

	rows, err := db.Query(sql)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var prescription = structs.Prescription{}

		err = rows.Scan(&prescription.ID, &prescription.Content, &prescription.PaymentLink, &prescription.PatientID, &prescription.DoctorID, &prescription.ConsultationID, &prescription.CreatedAt, &prescription.UpdatedAt)
		if err != nil {
			return
		}

		results = append(results, prescription)
	}

	return
}

func InsertPrescription(db *sql.DB, prescription structs.Prescription) (err error) {
	sql := "INSERT INTO prescription(id, content, payment_link, patient_id, doctor_id, consultation_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	errs := db.QueryRow(sql, prescription.ID, prescription.Content, prescription.PaymentLink, prescription.PatientID, prescription.DoctorID, prescription.ConsultationID, prescription.CreatedAt, prescription.UpdatedAt)

	return errs.Err()
}

func UpdatePrescription(db *sql.DB, prescription structs.Prescription) (err error) {

	sql := "UPDATE prescription SET content = $1, payment_link = $2, patient_id = $3, doctor_id = $4, consultation_id = $5, created_at = $6, updated_at = $7 WHERE id = $8"

	errs := db.QueryRow(sql, prescription.Content, prescription.PaymentLink, prescription.PatientID, prescription.DoctorID, prescription.ConsultationID, prescription.CreatedAt, prescription.UpdatedAt, prescription.ID)

	return errs.Err()
}

func DeletePrescription(db *sql.DB, prescription structs.Prescription) (err error) {
	sql := "DELETE FROM prescription WHERE id = $1"

	errs := db.QueryRow(sql, prescription.ID)

	return errs.Err()
}

func GetAllPrescriptionsByPatientID(db *sql.DB, patient structs.Patient) (err error, results []structs.Prescription) {
	sql := "SELECT * FROM prescription WHERE patient_id = $1"

	rows, err := db.Query(sql, patient.ID)
	if err != nil {
		return err, nil
	}

	defer rows.Close()

	for rows.Next() {
		var prescription = structs.Prescription{}

		err = rows.Scan(&prescription.ID, &prescription.Content, &prescription.PaymentLink, &prescription.PatientID, &prescription.DoctorID, &prescription.ConsultationID, &prescription.CreatedAt, &prescription.UpdatedAt)
		if err != nil {
			return err, nil
		}

		results = append(results, prescription)
	}

	return
}

func GetAllPrescriptionsByDoctorID(db *sql.DB, doctor structs.Doctor) (err error, results []structs.Prescription) {
	sql := "SELECT * FROM prescription WHERE doctor_id = $1"

	rows, err := db.Query(sql, doctor.ID)
	if err != nil {
		return err, nil
	}

	defer rows.Close()

	for rows.Next() {
		var prescription = structs.Prescription{}

		err = rows.Scan(&prescription.ID, &prescription.Content, &prescription.PaymentLink, &prescription.PatientID, &prescription.DoctorID, &prescription.ConsultationID, &prescription.CreatedAt, &prescription.UpdatedAt)
		if err != nil {
			return err, nil
		}

		results = append(results, prescription)
	}

	return
}

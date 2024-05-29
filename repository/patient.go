package repository

import (
	"database/sql"
	"telemedicine/structs"
)

func GetAllPatients(db *sql.DB) (err error, results []structs.Patient) {
	sql := "SELECT * FROM patient"

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var patient = structs.Patient{}

		err = rows.Scan(&patient.ID, &patient.Name, &patient.Password, &patient.CreatedAt, &patient.UpdatedAt)
		if err != nil {
			panic(err)
		}

		results = append(results, patient)
	}

	return
}

func InsertPatient(db *sql.DB, patient structs.Patient) (err error) {
	sql := "INSERT INTO patient(id, name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

	errs := db.QueryRow(sql, patient.ID, patient.Name, patient.Password, patient.CreatedAt, patient.UpdatedAt)

	return errs.Err()
}

func UpdatePatient(db *sql.DB, patient structs.Patient) (err error) {

	sql := "UPDATE patient SET name = $1, password = $2, created_at = $3, updated_at = $4 WHERE id = $5"

	errs := db.QueryRow(sql, patient.Name, patient.Password, patient.CreatedAt, patient.UpdatedAt, patient.ID)

	return errs.Err()
}

func DeletePatient(db *sql.DB, patient structs.Patient) (err error) {
	sql := "DELETE FROM patient WHERE id = $1"

	errs := db.QueryRow(sql, patient.ID)

	return errs.Err()
}

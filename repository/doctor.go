package repository

import (
	"database/sql"
	"telemedicine/structs"
)

func GetAllDoctors(db *sql.DB) (err error, results []structs.Doctor) {
	sql := "SELECT * FROM doctor"

	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var doctor = structs.Doctor{}

		err = rows.Scan(&doctor.ID, &doctor.Name, &doctor.Password, &doctor.CreatedAt, &doctor.UpdatedAt)
		if err != nil {
			panic(err)
		}

		results = append(results, doctor)
	}

	return
}

func InsertDoctor(db *sql.DB, doctor structs.Doctor) (err error) {
	sql := "INSERT INTO doctor(id, name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

	errs := db.QueryRow(sql, doctor.ID, doctor.Name, doctor.Password, doctor.CreatedAt, doctor.UpdatedAt)

	return errs.Err()
}

func UpdateDoctor(db *sql.DB, doctor structs.Doctor) (err error) {

	sql := "UPDATE doctor SET name = $1, password = $2, created_at = $3, updated_at = $4 WHERE id = $5"

	errs := db.QueryRow(sql, doctor.Name, doctor.Password, doctor.CreatedAt, doctor.UpdatedAt, doctor.ID)

	return errs.Err()
}

func DeleteDoctor(db *sql.DB, doctor structs.Doctor) (err error) {
	sql := "DELETE FROM doctor WHERE id = $1"

	errs := db.QueryRow(sql, doctor.ID)

	return errs.Err()
}

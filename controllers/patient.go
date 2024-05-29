package controllers

import (
	"net/http"
	"strconv"
	"telemedicine/database"
	"telemedicine/repository"
	"telemedicine/structs"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllPatients(c *gin.Context) {
	var (
		result gin.H
	)

	err, patients := repository.GetAllPatients(database.DbConnection)

	if err != nil {
		result = gin.H{
			"result": err,
		}
	} else {
		result = gin.H{
			"result": patients,
		}
	}

	c.JSON(http.StatusOK, result)
}

func InsertPatient(c *gin.Context) {
	var patient structs.Patient

	err := c.ShouldBindJSON(&patient)
	if err != nil {
		panic(err)
	}

	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()

	err2, patients := repository.GetAllPatients(database.DbConnection)
	patient.ID = 0
	if err2 != nil {
		panic(err2)
	}
	for _, v := range patients {
		if v.ID > patient.ID {
			// Mengambil nilai maksimum indeks ID
			patient.ID = v.ID
		}
	}
	patient.ID++

	err = repository.InsertPatient(database.DbConnection, patient)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Success Insert Patient",
	})

}

func UpdatePatient(c *gin.Context) {
	var patient structs.Patient
	id, _ := strconv.Atoi(c.Param("id"))

	err := c.ShouldBindJSON(&patient)
	if err != nil {
		panic(err)
	}

	patient.ID = int64(id)
	patient.UpdatedAt = time.Now()

	err2, patients := repository.GetAllPatients(database.DbConnection)
	if err2 != nil {
		panic(err2)
	}
	for _, v := range patients {
		if v.ID == patient.ID {
			// Mengambil waktu CreatedAt pada database sebelum diupdate
			patient.CreatedAt = v.CreatedAt
		}
	}

	// Update Patient
	err = repository.UpdatePatient(database.DbConnection, patient)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Success Update Patient",
	})
}

func DeletePatient(c *gin.Context) {
	var patient structs.Patient
	id, _ := strconv.Atoi(c.Param("id"))

	patient.ID = int64(id)

	err := repository.DeletePatient(database.DbConnection, patient)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Success Delete Patient",
	})
}

package controllers

import (
	"net/http"
	"strconv"
	"telemedicine/database"
	"telemedicine/repository"
	"telemedicine/structs"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllDoctors(c *gin.Context) {
	var (
		result gin.H
	)

	err, doctors := repository.GetAllDoctors(database.DbConnection)

	if err != nil {
		result = gin.H{
			"result": err,
		}
	} else {
		result = gin.H{
			"result": doctors,
		}
	}

	c.JSON(http.StatusOK, result)
}

func InsertDoctor(c *gin.Context) {
	var doctor structs.Doctor

	err := c.ShouldBindJSON(&doctor)
	if err != nil {
		panic(err)
	}

	doctor.CreatedAt = time.Now()
	doctor.UpdatedAt = time.Now()

	err2, doctors := repository.GetAllDoctors(database.DbConnection)
	doctor.ID = 0
	if err2 != nil {
		panic(err2)
	}
	for _, v := range doctors {
		if v.ID > doctor.ID {
			// Mengambil nilai maksimum indeks ID
			doctor.ID = v.ID
		}
	}
	doctor.ID++

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	doctor.Password = string(hashedPassword)

	err = repository.InsertDoctor(database.DbConnection, doctor)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Success Insert Doctor",
	})

}

func UpdateDoctor(c *gin.Context) {
	var doctor structs.Doctor
	id, _ := strconv.Atoi(c.Param("id"))

	err := c.ShouldBindJSON(&doctor)
	if err != nil {
		panic(err)
	}

	doctor.ID = int64(id)
	doctor.UpdatedAt = time.Now()

	err2, doctors := repository.GetAllDoctors(database.DbConnection)
	if err2 != nil {
		panic(err2)
	}
	for _, v := range doctors {
		if v.ID == doctor.ID {
			// Mengambil waktu CreatedAt pada database sebelum diupdate
			doctor.CreatedAt = v.CreatedAt
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	doctor.Password = string(hashedPassword)

	// Update Doctor
	err = repository.UpdateDoctor(database.DbConnection, doctor)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Success Update Doctor",
	})
}

func DeleteDoctor(c *gin.Context) {
	var doctor structs.Doctor
	id, _ := strconv.Atoi(c.Param("id"))

	doctor.ID = int64(id)

	err := repository.DeleteDoctor(database.DbConnection, doctor)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "Success Delete Doctor",
	})
}

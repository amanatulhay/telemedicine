package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"telemedicine/database"
	"telemedicine/repository"
	"telemedicine/structs"
	"telemedicine/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllPatients(c *gin.Context) {

	err, patients := repository.GetAllPatients(database.DbConnection)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Berhasil mengambil seluruh data Patients",
			"data":    patients,
		})
	}
}

func IsDuplicatePatient(name string) (bool, error) {

	err, patients := repository.GetAllPatients(database.DbConnection)
	if err != nil {
		return true, err
	}

	for _, v := range patients {
		if v.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func InsertPatient(c *gin.Context) {
	var patient structs.Patient

	err := c.ShouldBindJSON(&patient)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Field Name dan Password tidak boleh kosong",
			"data":    utils.NullData,
		})
		return
	}

	isDuplicate, err := IsDuplicatePatient(patient.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	if isDuplicate {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan Name %s telah disimpan", patient.Name),
			"data":    utils.NullData,
		})
		return
	}

	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()

	err2, patients := repository.GetAllPatients(database.DbConnection)
	patient.ID = 0
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	for _, v := range patients {
		if v.ID > patient.ID {
			// Mengambil nilai maksimum indeks ID
			patient.ID = v.ID
		}
	}
	patient.ID++

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	patient.Password = string(hashedPassword)

	err = repository.InsertPatient(database.DbConnection, patient)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil menambahkan data Patient",
		"data":    utils.NullData,
	})
}

func UpdatePatient(c *gin.Context) {
	var patient structs.Patient
	var createdAt time.Time
	id, _ := strconv.Atoi(c.Param("id"))

	err, patients := repository.GetAllPatients(database.DbConnection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	found := false
	for _, v := range patients {
		if v.ID == int64(id) {
			// Mengambil waktu CreatedAt pada database sebelum diupdate
			createdAt = v.CreatedAt
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan id %d tidak ditemukan", id),
			"data":    utils.NullData,
		})
		return
	}

	err = c.ShouldBindJSON(&patient)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Field Name dan Password tidak boleh kosong",
			"data":    utils.NullData,
		})
		return
	}

	isDuplicate, err := IsDuplicatePatient(patient.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	if isDuplicate {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan Name %s telah disimpan", patient.Name),
			"data":    utils.NullData,
		})
		return
	}

	patient.ID = int64(id)
	patient.UpdatedAt = time.Now()
	patient.CreatedAt = createdAt

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	patient.Password = string(hashedPassword)

	// Update Patient
	err = repository.UpdatePatient(database.DbConnection, patient)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil memperbarui data Patient",
		"data":    utils.NullData,
	})
}

func DeletePatient(c *gin.Context) {
	var patient structs.Patient
	id, _ := strconv.Atoi(c.Param("id"))

	err, patients := repository.GetAllPatients(database.DbConnection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	found := false
	for _, v := range patients {
		if v.ID == int64(id) {
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan id %d tidak ditemukan", id),
			"data":    utils.NullData,
		})
		return
	}

	patient.ID = int64(id)

	err = repository.DeletePatient(database.DbConnection, patient)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil menghapus data Patient",
		"data":    utils.NullData,
	})
}

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

func GetAllDoctors(c *gin.Context) {

	err, doctors := repository.GetAllDoctors(database.DbConnection)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Berhasil mengambil seluruh data Doctors",
			"data":    doctors,
		})
	}
}

func IsDuplicateDoctor(name string) (bool, error) {

	err, doctors := repository.GetAllDoctors(database.DbConnection)
	if err != nil {
		return true, err
	}

	for _, v := range doctors {
		if v.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func InsertDoctor(c *gin.Context) {
	var doctor structs.Doctor

	err := c.ShouldBindJSON(&doctor)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Field Name dan Password tidak boleh kosong",
			"data":    utils.NullData,
		})
		return
	}

	isDuplicate, err := IsDuplicateDoctor(doctor.Name)
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
			"message": fmt.Sprintf("Data Doctor dengan Name %s telah disimpan", doctor.Name),
			"data":    utils.NullData,
		})
		return
	}

	doctor.CreatedAt = time.Now()
	doctor.UpdatedAt = time.Now()

	err2, doctors := repository.GetAllDoctors(database.DbConnection)
	doctor.ID = 0
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	doctor.Password = string(hashedPassword)

	err = repository.InsertDoctor(database.DbConnection, doctor)
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
		"message": "Berhasil menambahkan data Doctor",
		"data":    utils.NullData,
	})
}

func UpdateDoctor(c *gin.Context) {
	var doctor structs.Doctor
	var createdAt time.Time
	id, _ := strconv.Atoi(c.Param("id"))

	err, doctors := repository.GetAllDoctors(database.DbConnection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	found := false
	for _, v := range doctors {
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
			"message": fmt.Sprintf("Data Doctor dengan id %d tidak ditemukan", id),
			"data":    utils.NullData,
		})
		return
	}

	err = c.ShouldBindJSON(&doctor)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Field Name dan Password tidak boleh kosong",
			"data":    utils.NullData,
		})
		return
	}

	isDuplicate, err := IsDuplicateDoctor(doctor.Name)
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
			"message": fmt.Sprintf("Data Doctor dengan Name %s telah disimpan", doctor.Name),
			"data":    utils.NullData,
		})
		return
	}

	doctor.ID = int64(id)
	doctor.UpdatedAt = time.Now()
	doctor.CreatedAt = createdAt

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(doctor.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	doctor.Password = string(hashedPassword)

	// Update Doctor
	err = repository.UpdateDoctor(database.DbConnection, doctor)

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
		"message": "Berhasil memperbarui data Doctor",
		"data":    utils.NullData,
	})
}

func DeleteDoctor(c *gin.Context) {
	var doctor structs.Doctor
	id, _ := strconv.Atoi(c.Param("id"))

	err, doctors := repository.GetAllDoctors(database.DbConnection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	found := false
	for _, v := range doctors {
		if v.ID == int64(id) {
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Doctor dengan id %d tidak ditemukan", id),
			"data":    utils.NullData,
		})
		return
	}

	doctor.ID = int64(id)

	err = repository.DeleteDoctor(database.DbConnection, doctor)
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
		"message": "Berhasil menghapus data Doctor",
		"data":    utils.NullData,
	})
}

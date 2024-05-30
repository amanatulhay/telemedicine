package controllers

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"telemedicine/database"
	"telemedicine/repository"
	"telemedicine/structs"
	"telemedicine/utils"
	"time"

	"telemedicine/utils/token"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CurrentPatient(c *gin.Context) {

	var p structs.Patient

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    utils.NullData,
		})
		return
	}

	// Mencari data user berdasarkan ID
	err2, patients := repository.GetAllPatients(database.DbConnection)
	found := false
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	for _, v := range patients {
		if user_id == uint(v.ID) {
			p = v
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan id %d tidak ditemukan", user_id),
			"data":    utils.NullData,
		})
		return
	}

	p.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail data Patient",
		"data":    p,
	})
}

func LoginPatient(c *gin.Context) {

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

	err2, patients := repository.GetAllPatients(database.DbConnection)
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	var hashedPassword string
	var ID uint
	for _, v := range patients {
		if v.Name == patient.Name {
			// Mengambil nilai hashedPassword pada database sesuai Nama pada login
			hashedPassword = v.Password
			ID = uint(v.ID)
			break
		}
	}
	if hashedPassword == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan Name %s tidak ditemukan", patient.Name),
			"data":    utils.NullData,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(patient.Password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Password is incorrect",
			"data":    utils.NullData,
		})
		return
	}

	token, err := token.GenerateToken(ID)

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
		"message": "Berhasil login akun Patient",
		"token":   token,
	})

}

func RegisterPatient(c *gin.Context) {

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

	//remove spaces in name
	patient.Name = html.EscapeString(strings.TrimSpace(patient.Name))

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

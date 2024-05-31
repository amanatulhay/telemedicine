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

func CurrentDoctorData(c *gin.Context) {

	var d structs.Doctor

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
	err2, doctors := repository.GetAllDoctors(database.DbConnection)
	found := false
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	for _, v := range doctors {
		if user_id == uint(v.ID) {
			d = v
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Doctor dengan id %d tidak ditemukan", user_id),
			"data":    utils.NullData,
		})
		return
	}

	d.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail data Doctor",
		"data":    d,
	})
}

func CurrentDoctorConsultations(c *gin.Context) {

	var d structs.Doctor

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
	err2, doctors := repository.GetAllDoctors(database.DbConnection)
	found := false
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	for _, v := range doctors {
		if user_id == uint(v.ID) {
			d = v
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Doctor dengan id %d tidak ditemukan", user_id),
			"data":    utils.NullData,
		})
		return
	}

	err, consultations := repository.GetAllConsultationsByDoctorID(database.DbConnection, d)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Berhasil mengambil seluruh data Consultations berdasarkan Token ID Doctor %d", d.ID),
			"data":    consultations,
		})
	}
}

func LoginDoctor(c *gin.Context) {

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

	err2, doctors := repository.GetAllDoctors(database.DbConnection)
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
	for _, v := range doctors {
		if v.Name == doctor.Name {
			// Mengambil nilai hashedPassword pada database sesuai Nama pada login
			hashedPassword = v.Password
			ID = uint(v.ID)
			break
		}
	}
	if hashedPassword == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Doctor dengan Name %s tidak ditemukan", doctor.Name),
			"data":    utils.NullData,
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(doctor.Password))

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
		"message": "Berhasil login akun Doctor",
		"token":   token,
	})

}

func RegisterDoctor(c *gin.Context) {

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

	//remove spaces in name
	doctor.Name = html.EscapeString(strings.TrimSpace(doctor.Name))

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

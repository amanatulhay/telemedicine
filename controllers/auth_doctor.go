package controllers

import (
	"errors"
	"html"
	"net/http"
	"strings"
	"telemedicine/database"
	"telemedicine/repository"
	"telemedicine/structs"
	"time"

	"telemedicine/utils/token"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CurrentDoctor(c *gin.Context) {

	var p structs.Doctor

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mencari data user berdasarkan ID
	err2, doctors := repository.GetAllDoctors(database.DbConnection)
	found := false
	if err2 != nil {
		panic(err2)
	}
	for _, v := range doctors {
		if user_id == uint(v.ID) {
			p = v
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Doctor not found!")})
		return
	}

	p.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": p})
}

func LoginDoctor(c *gin.Context) {

	var doctor structs.Doctor

	err := c.ShouldBindJSON(&doctor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err2, doctors := repository.GetAllDoctors(database.DbConnection)
	if err2 != nil {
		panic(err2)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(doctor.Password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name or password is incorrect."})
		return
	}

	token, err := token.GenerateToken(ID)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func RegisterDoctor(c *gin.Context) {

	var doctor structs.Doctor

	err := c.ShouldBindJSON(&doctor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	//remove spaces in name
	doctor.Name = html.EscapeString(strings.TrimSpace(doctor.Name))

	err = repository.InsertDoctor(database.DbConnection, doctor)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Doctor registration success"})

}

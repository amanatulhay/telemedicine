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

func CurrentUser(c *gin.Context) {

	var p structs.Patient

	user_id, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mencari data user berdasarkan ID
	err2, patients := repository.GetAllPatients(database.DbConnection)
	found := false
	if err2 != nil {
		panic(err2)
	}
	for _, v := range patients {
		if user_id == uint(v.ID) {
			p = v
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("User not found!")})
		return
	}

	p.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": p})
}

func Login(c *gin.Context) {

	var patient structs.Patient

	err := c.ShouldBindJSON(&patient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// u := models.User{}

	// u.Username = input.Username
	// u.Password = input.Password

	// token, err := models.LoginCheck(u.Username, u.Password)

	err2, patients := repository.GetAllPatients(database.DbConnection)
	if err2 != nil {
		panic(err2)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(patient.Password))

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

func Register(c *gin.Context) {

	var patient structs.Patient

	err := c.ShouldBindJSON(&patient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	patient.Password = string(hashedPassword)

	//remove spaces in username
	patient.Name = html.EscapeString(strings.TrimSpace(patient.Name))

	err = repository.InsertPatient(database.DbConnection, patient)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})

}

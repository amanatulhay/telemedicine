package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"telemedicine/database"
	"telemedicine/repository"
	"telemedicine/structs"
	"telemedicine/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllConsultations(c *gin.Context) {

	err, consultations := repository.GetAllConsultations(database.DbConnection)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Berhasil mengambil seluruh data Consultations",
			"data":    consultations,
		})
	}
}

func InsertConsultation(c *gin.Context) {
	var consultation structs.Consultation

	err := c.ShouldBindJSON(&consultation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Field meeting_link, payment_link, patient_id, dan doctor_id tidak boleh kosong",
			"data":    utils.NullData,
		})
		return
	}

	match1, _ := regexp.MatchString("https?://.*", consultation.MeetingLink)
	match2, _ := regexp.MatchString("https?://.*", consultation.PaymentLink)
	if !match1 && !match2 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success":   false,
			"message 1": "meeting_link's format should be an URL (e.g. : https://... or http://...) ",
			"message 2": "payment_link's format should be an URL (e.g. : https://... or http://...) ",
			"data":      utils.NullData,
		})
		return
	} else if !match1 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "meeting_link's format should be an URL (e.g. : https://... or http://...) ",
			"data":    utils.NullData,
		})
		return
	} else if !match2 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "payment_link's format should be an URL (e.g. : https://... or http://...) ",
			"data":    utils.NullData,
		})
		return
	}

	idExists, err := PatientIDExists(consultation.PatientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	if !idExists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan ID %d tidak ada", consultation.PatientID),
			"data":    utils.NullData,
		})
		return
	}

	idExists, err = DoctorIDExists(consultation.DoctorID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	if !idExists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Doctor dengan ID %d tidak ada", consultation.DoctorID),
			"data":    utils.NullData,
		})
		return
	}

	consultation.CreatedAt = time.Now()
	consultation.UpdatedAt = time.Now()

	err2, consultations := repository.GetAllConsultations(database.DbConnection)
	consultation.ID = 0
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	for _, v := range consultations {
		if v.ID > consultation.ID {
			// Mengambil nilai maksimum indeks ID
			consultation.ID = v.ID
		}
	}
	consultation.ID++

	err = repository.InsertConsultation(database.DbConnection, consultation)
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
		"message": "Berhasil menambahkan data Consultation",
		"data":    utils.NullData,
	})
}

func UpdateConsultation(c *gin.Context) {
	var consultation structs.Consultation
	var createdAt time.Time
	id, _ := strconv.Atoi(c.Param("id"))

	err, consultations := repository.GetAllConsultations(database.DbConnection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	found := false
	for _, v := range consultations {
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
			"message": fmt.Sprintf("Data Consultation dengan id %d tidak ditemukan", id),
			"data":    utils.NullData,
		})
		return
	}

	err = c.ShouldBindJSON(&consultation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Field meeting_link, payment_link, patient_id, dan doctor_id tidak boleh kosong",
			"data":    utils.NullData,
		})
		return
	}

	match1, _ := regexp.MatchString("https?://.*", consultation.MeetingLink)
	match2, _ := regexp.MatchString("https?://.*", consultation.PaymentLink)
	if !match1 && !match2 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success":   false,
			"message 1": "meeting_link's format should be an URL (e.g. : https://... or http://...) ",
			"message 2": "payment_link's format should be an URL (e.g. : https://... or http://...) ",
			"data":      utils.NullData,
		})
		return
	} else if !match1 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "meeting_link's format should be an URL (e.g. : https://... or http://...) ",
			"data":    utils.NullData,
		})
		return
	} else if !match2 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "payment_link's format should be an URL (e.g. : https://... or http://...) ",
			"data":    utils.NullData,
		})
		return
	}

	idExists, err := PatientIDExists(consultation.PatientID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	if !idExists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan ID %d tidak ada", consultation.PatientID),
			"data":    utils.NullData,
		})
		return
	}

	idExists, err = DoctorIDExists(consultation.DoctorID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	if !idExists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Doctor dengan ID %d tidak ada", consultation.DoctorID),
			"data":    utils.NullData,
		})
		return
	}

	consultation.ID = int64(id)
	consultation.UpdatedAt = time.Now()
	consultation.CreatedAt = createdAt

	// Update Consultation
	err = repository.UpdateConsultation(database.DbConnection, consultation)

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
		"message": "Berhasil memperbarui data Consultation",
		"data":    utils.NullData,
	})
}

func DeleteConsultation(c *gin.Context) {
	var consultation structs.Consultation
	id, _ := strconv.Atoi(c.Param("id"))

	err, consultations := repository.GetAllConsultations(database.DbConnection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	found := false
	for _, v := range consultations {
		if v.ID == int64(id) {
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Consultation dengan id %d tidak ditemukan", id),
			"data":    utils.NullData,
		})
		return
	}

	consultation.ID = int64(id)

	err = repository.DeleteConsultation(database.DbConnection, consultation)
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
		"message": "Berhasil menghapus data Consultation",
		"data":    utils.NullData,
	})
}

func GetAllConsultationsByPatientID(c *gin.Context) {

	var patient structs.Patient
	id, _ := strconv.Atoi(c.Param("id"))

	patient.ID = int64(id)

	idExists, err := PatientIDExists(patient.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	if !idExists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Patient dengan ID %d tidak ada", patient.ID),
			"data":    utils.NullData,
		})
		return
	}

	err, consultations := repository.GetAllConsultationsByPatientID(database.DbConnection, patient)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Berhasil mengambil seluruh data Consultations berdasarkan ID Pasien %d", patient.ID),
			"data":    consultations,
		})
	}
}

func PatientIDExists(id int64) (bool, error) {

	err, patients := repository.GetAllPatients(database.DbConnection)
	if err != nil {
		return false, err
	}

	for _, v := range patients {
		if v.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func GetAllConsultationsByDoctorID(c *gin.Context) {

	var doctor structs.Doctor
	id, _ := strconv.Atoi(c.Param("id"))

	doctor.ID = int64(id)

	idExists, err := DoctorIDExists(doctor.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	if !idExists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Doctor dengan ID %d tidak ada", doctor.ID),
			"data":    utils.NullData,
		})
		return
	}

	err, consultations := repository.GetAllConsultationsByDoctorID(database.DbConnection, doctor)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Berhasil mengambil seluruh data Consultations berdasarkan ID Doctor %d", doctor.ID),
			"data":    consultations,
		})
	}
}

func DoctorIDExists(id int64) (bool, error) {

	err, doctors := repository.GetAllDoctors(database.DbConnection)
	if err != nil {
		return false, err
	}

	for _, v := range doctors {
		if v.ID == id {
			return true, nil
		}
	}
	return false, nil
}

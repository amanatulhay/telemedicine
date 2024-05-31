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

func GetAllPrescriptions(c *gin.Context) {

	err, prescriptions := repository.GetAllPrescriptions(database.DbConnection)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Berhasil mengambil seluruh data Prescriptions",
			"data":    prescriptions,
		})
	}
}

func InsertPrescription(c *gin.Context) {
	var prescription structs.Prescription

	err := c.ShouldBindJSON(&prescription)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Field content, payment_link, patient_id, doctor_id, dan consultation_id tidak boleh kosong",
			"data":    utils.NullData,
		})
		return
	}

	match, _ := regexp.MatchString("https?://.*", prescription.PaymentLink)
	if !match {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "payment_link's format should be an URL (e.g. : https://... or http://...) ",
			"data":    utils.NullData,
		})
		return
	}

	idExists, err := PatientIDExists(prescription.PatientID)
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
			"message": fmt.Sprintf("Data Patient dengan ID %d tidak ada", prescription.PatientID),
			"data":    utils.NullData,
		})
		return
	}

	idExists, err = DoctorIDExists(prescription.DoctorID)
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
			"message": fmt.Sprintf("Data Doctor dengan ID %d tidak ada", prescription.DoctorID),
			"data":    utils.NullData,
		})
		return
	}

	idExists, err = ConsultationIDExists(prescription.ConsultationID)
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
			"message": fmt.Sprintf("Data Consultation dengan ID %d tidak ada", prescription.ConsultationID),
			"data":    utils.NullData,
		})
		return
	}

	prescription.CreatedAt = time.Now()
	prescription.UpdatedAt = time.Now()

	err2, prescriptions := repository.GetAllPrescriptions(database.DbConnection)
	prescription.ID = 0
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}
	for _, v := range prescriptions {
		if v.ID > prescription.ID {
			// Mengambil nilai maksimum indeks ID
			prescription.ID = v.ID
		}
	}
	prescription.ID++

	err = repository.InsertPrescription(database.DbConnection, prescription)
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
		"message": "Berhasil menambahkan data Prescription",
		"data":    utils.NullData,
	})
}

func UpdatePrescription(c *gin.Context) {
	var prescription structs.Prescription
	var createdAt time.Time
	id, _ := strconv.Atoi(c.Param("id"))

	err, prescriptions := repository.GetAllPrescriptions(database.DbConnection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	found := false
	for _, v := range prescriptions {
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
			"message": fmt.Sprintf("Data Prescription dengan id %d tidak ditemukan", id),
			"data":    utils.NullData,
		})
		return
	}

	err = c.ShouldBindJSON(&prescription)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Field content, payment_link, patient_id, doctor_id, dan consultation_id tidak boleh kosong",
			"data":    utils.NullData,
		})
		return
	}

	match, _ := regexp.MatchString("https?://.*", prescription.PaymentLink)
	if !match {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "payment_link's format should be an URL (e.g. : https://... or http://...) ",
			"data":    utils.NullData,
		})
		return
	}

	idExists, err := PatientIDExists(prescription.PatientID)
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
			"message": fmt.Sprintf("Data Patient dengan ID %d tidak ada", prescription.PatientID),
			"data":    utils.NullData,
		})
		return
	}

	idExists, err = DoctorIDExists(prescription.DoctorID)
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
			"message": fmt.Sprintf("Data Doctor dengan ID %d tidak ada", prescription.DoctorID),
			"data":    utils.NullData,
		})
		return
	}

	idExists, err = ConsultationIDExists(prescription.ConsultationID)
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
			"message": fmt.Sprintf("Data Consultation dengan ID %d tidak ada", prescription.ConsultationID),
			"data":    utils.NullData,
		})
		return
	}

	prescription.ID = int64(id)
	prescription.UpdatedAt = time.Now()
	prescription.CreatedAt = createdAt

	// Update Prescription
	err = repository.UpdatePrescription(database.DbConnection, prescription)

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
		"message": "Berhasil memperbarui data Prescription",
		"data":    utils.NullData,
	})
}

func DeletePrescription(c *gin.Context) {
	var prescription structs.Prescription
	id, _ := strconv.Atoi(c.Param("id"))

	err, prescriptions := repository.GetAllPrescriptions(database.DbConnection)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
		return
	}

	found := false
	for _, v := range prescriptions {
		if v.ID == int64(id) {
			found = true
			break
		}
	}

	if !found {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Data Prescription dengan id %d tidak ditemukan", id),
			"data":    utils.NullData,
		})
		return
	}

	prescription.ID = int64(id)

	err = repository.DeletePrescription(database.DbConnection, prescription)
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
		"message": "Berhasil menghapus data Prescription",
		"data":    utils.NullData,
	})
}

func GetAllPrescriptionsByPatientID(c *gin.Context) {

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

	err, prescriptions := repository.GetAllPrescriptionsByPatientID(database.DbConnection, patient)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Berhasil mengambil seluruh data Prescriptions berdasarkan ID Pasien %d", patient.ID),
			"data":    prescriptions,
		})
	}
}

func GetAllPrescriptionsByDoctorID(c *gin.Context) {

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

	err, prescriptions := repository.GetAllPrescriptionsByDoctorID(database.DbConnection, doctor)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
			"data":    utils.NullData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Berhasil mengambil seluruh data Prescriptions berdasarkan ID Doctor %d", doctor.ID),
			"data":    prescriptions,
		})
	}
}

func ConsultationIDExists(id int64) (bool, error) {

	err, consultations := repository.GetAllConsultations(database.DbConnection)
	if err != nil {
		return false, err
	}

	for _, v := range consultations {
		if v.ID == id {
			return true, nil
		}
	}
	return false, nil
}

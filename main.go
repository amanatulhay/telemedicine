package main

import (
	"database/sql"
	"fmt"
	"os"
	"telemedicine/controllers"
	"telemedicine/database"
	"telemedicine/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func main() {

	err = godotenv.Load("config/.env")
	if err != nil {
		fmt.Println("failed load file environment")
	} else {
		fmt.Println("success read file environment")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	DB, err = sql.Open("postgres", psqlInfo)
	err = DB.Ping()
	if err != nil {
		fmt.Println("DB Connection Failed")
		panic(err)
	} else {
		fmt.Println("DB Connection Success")
	}

	database.DbMigrate(DB)

	defer DB.Close()

	// Router GIN
	router := gin.Default()

	// Router Admin
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin":  "password",
		"editor": "secret",
	}))

	// Router Admin - Consultations
	authorized.GET("/consultations", controllers.GetAllConsultations)
	authorized.POST("/consultations", controllers.InsertConsultation)
	authorized.PUT("/consultations/:id", controllers.UpdateConsultation)
	authorized.DELETE("/consultations/:id", controllers.DeleteConsultation)
	authorized.GET("/patients/:id/consultations", controllers.GetAllConsultationsByPatientID)
	authorized.GET("/doctors/:id/consultations", controllers.GetAllConsultationsByDoctorID)

	// Router Admin - Doctors
	authorized.GET("/doctors", controllers.GetAllDoctors)
	authorized.POST("/doctors", controllers.InsertDoctor)
	authorized.PUT("/doctors/:id", controllers.UpdateDoctor)
	authorized.DELETE("/doctors/:id", controllers.DeleteDoctor)

	// Router Admin - Patients
	authorized.GET("/patients", controllers.GetAllPatients)
	authorized.POST("/patients", controllers.InsertPatient)
	authorized.PUT("/patients/:id", controllers.UpdatePatient)
	authorized.DELETE("/patients/:id", controllers.DeletePatient)

	// Router Admin - Prescriptions
	authorized.GET("/prescriptions", controllers.GetAllPrescriptions)
	authorized.POST("/prescriptions", controllers.InsertPrescription)
	authorized.PUT("/prescriptions/:id", controllers.UpdatePrescription)
	authorized.DELETE("/prescriptions/:id", controllers.DeletePrescription)
	authorized.GET("/patients/:id/prescriptions", controllers.GetAllPrescriptionsByPatientID)
	authorized.GET("/doctors/:id/prescriptions", controllers.GetAllPrescriptionsByDoctorID)

	// Router JWT Auth
	protected := router.Group("/api/")
	protected.Use(middlewares.JwtAuthMiddleware())

	protected.POST("/current-patient-add-consultation", controllers.CurrentPatientAddConsultation)
	protected.GET("/current-patient-data", controllers.CurrentPatientData)
	protected.GET("/current-patient-all-consultations", controllers.CurrentPatientAllConsultations)
	protected.GET("/current-patient-all-prescriptions", controllers.CurrentPatientAllPrescriptions)

	protected.POST("/current-doctor-add-prescription", controllers.CurrentDoctorAddPrescription)
	protected.GET("/current-doctor-data", controllers.CurrentDoctorData)
	protected.GET("/current-doctor-all-consultations", controllers.CurrentDoctorAllConsultations)
	protected.GET("/current-doctor-all-prescriptions", controllers.CurrentDoctorAllPrescriptions)

	// Router Public Doctor
	publicDoctor := router.Group("/doctor")

	publicDoctor.POST("/register", controllers.RegisterDoctor)
	publicDoctor.POST("/login", controllers.LoginDoctor)

	// Router Public Patient
	publicPatient := router.Group("/patient")

	publicPatient.POST("/register", controllers.RegisterPatient)
	publicPatient.POST("/login", controllers.LoginPatient)

	router.Run("localhost:8080")
}

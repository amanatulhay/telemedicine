package main

import (
	"database/sql"
	"fmt"
	"telemedicine/controllers"
	"telemedicine/database"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "telemedicine"
)

var (
	DB  *sql.DB
	err error
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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

	// Router Category
	router.GET("/patients", controllers.GetAllPatients)
	router.POST("/patients", controllers.InsertPatient)
	router.PUT("/patients/:id", controllers.UpdatePatient)
	router.DELETE("/patients/:id", controllers.DeletePatient)

	router.Run("localhost:8080")
}

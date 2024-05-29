package main

import (
	"database/sql"
	"fmt"
	"telemedicine/controllers"
	"telemedicine/database"
	"telemedicine/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("failed load file environment")
	} else {
		fmt.Println("success read file environment")
	}

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

	// Router Admin
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin":  "password",
		"editor": "secret",
	}))

	authorized.GET("/patients", controllers.GetAllPatients)
	authorized.DELETE("/patients/:id", controllers.DeletePatient)

	// Router Public
	public := router.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	// Router JWT Auth
	protected := router.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())

	protected.GET("/user", controllers.CurrentUser)

	router.Run("localhost:8080")
}

package main

import (
	"fmt"
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	host := "postgres"
	port := "5432"
	user := "postgres"
	password := "postgres"
	dbname := "app"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return &Database{DB: db}
}

func main() {
    r := gin.Default()
	db := NewDatabase()
	defer db.DB.Close()

	r.POST("/signup", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	r.Run(":3000")
}
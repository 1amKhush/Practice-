package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	DbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	
	var err error
	DB, err = sql.Open("postgres", DbInfo)

	if err != nil {
		log.Fatal("Error opening DB : ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error connecting to DB : ", err)
	}

	log.Println()
	fmt.Println("-> Succesfully connected to DB")
}
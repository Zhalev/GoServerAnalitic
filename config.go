package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const (
	host     = "192.168.1.45"
	port     = 5433
	user     = "postgres"
	password = "123"
	dbname   = "devices_management"
)

func Config() {

	f, err1 := os.OpenFile("log_backend.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	defer f.Close()

	logger := log.New(f, "ERROR:", log.LstdFlags)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	defer db.Close()

	if err = db.Ping(); err != nil {
		logger.Panic("Failed to get config: ", err)
	} else {
		logger.Println("DB connected successfully")
	}

}

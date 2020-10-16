package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	_ "time"
	/*	_ "db.properties"*/)

const (
	host     = "192.168.1.45"
	port     = 5433
	user     = "postgres"
	password = "123"
	dbname   = "devices_management"
)

func main() {

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
	if err != nil {
		logger.Fatal("Connect crash:", err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		logger.Panic("Failed to get config: ", err)
		db.Close()
	} else {
		logger.Println("DB connected successfully")
	}

	var (
		id                  int
		name                string
		deviceStatus        interface{}
		placeOfInstallation interface{}
		yearOfInstallation  interface{}
		written             interface{}
	)

	rows, err := db.Query("select id, name, device_status, place_of_installation, year_of_installation, written from  devices  order by id")
	if err != nil {
		logger.Fatal("SQL:", err, rows)
	}
	defer rows.Close()
	fmt.Println("id |   name    | deviceStatus | placeOfInstallation | yearOfInstallation | written")
	for rows.Next() {
		err := rows.Scan(&id, &name, &deviceStatus, &placeOfInstallation, &yearOfInstallation, &written)
		if err != nil {
			logger.Fatal(err)
		} else {
			logger.Println("Successfully2")
		}
		fmt.Printf("%2v | %9v | %12v | %19v | %18v | %6v\n", id, name, deviceStatus, placeOfInstallation, yearOfInstallation, written)
	}

	err = rows.Err()
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Println("Successfully3")
	}

}

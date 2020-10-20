package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"./config"

)

var (
	name        interface{}
	startDate	interface{}
	endDate		interface{}
)

func main() {

	http.HandleFunc("/plan_log",CreatePlanLog)
	http.HandleFunc("/plan_log/{id}",UpdatePlanLog)
	http.HandleFunc("/plan_log",SelectPlanLog)
	
}

func CreatePlanLog(w http.ResponseWriter, r *http.Request) {

}

func UpdatePlanLog(w http.ResponseWriter, r *http.Request)  {

}

func SelectPlanLog(w http.ResponseWriter, r *http.Request) {
	f, err1 := os.OpenFile("log_backend.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		log.Println(err1)
	}
	defer f.Close()

	logger := log.New(f, "ERROR:", log.LstdFlags)
	// import config
	db, err := sql.Open("postgres",  psqlInfo)

	rows, err := db.Query("select * from plan_log")
	if err != nil {
		logger.Fatal("SQL:", err, rows)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name, &startDate, &endDate)
		if err != nil {
			logger.Fatal(err)
		}
	}
}

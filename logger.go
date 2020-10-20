package main

import (
	"log"
	"os"
)

func Logger()  {

	f, err := os.OpenFile("log_backend.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
}

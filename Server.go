package main

import (
	"fmt"
	"log"
	"net/http"
)

type Handler struct {
	Name string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Name:", h.Name, "URL:", r.URL.String())
}
func main() {
	// можно из метода структуры обрабочик делать
	//testHandler := &Handler{Name: "test"}
	//http.Handle("/test/", testHandler)
	//rootHandler := &Handler{Name: "root"}
	//http.Handle("/", rootHandler)
	// а можно из функции
	http.HandleFunc("/", handler)
	fmt.Println("starting server at :8080")
	//http.ListenAndServe(":8080", nil)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// чтение параметров или так
	myParam := r.URL.Query().Get("param")
	if myParam != "" {
		fmt.Fprintln(w, "‘myParam‘ is", myParam)
	}
	// или так так получаем как get, так и post параметры
	key := r.FormValue("key")
	if key != "" {
		fmt.Fprintln(w, "‘key‘ is", key)
	}
	if key == "1" {
		fmt.Fprintln(w, "You browser is", r.UserAgent())
		http.Redirect(w, r, "/", http.StatusFound)
	}
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Fprintf(w, "yasdkf")
}

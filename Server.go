package main

import (
	"fmt"
	"log"
	"net"
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
	fmt.Println("starting server at :5555")
	//http.ListenAndServe(":8080", nil)
	log.Fatal(http.ListenAndServe(":5555", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// чтение параметров или так
	myParam := r.URL.Query().Get("param")
	//w.Header().Set("Content-Type","text/html")
	if myParam != "" {
		fmt.Fprintln(w, "‘myParam‘ is", myParam)
	}
	// или так так получаем как get, так и post параметры
	name := r.FormValue("name")
	if name != "" {
		fmt.Fprintln(w, "Hello, there is", name)
	}
	if name == "1" {
		fmt.Fprintln(w, "You browser is", r.UserAgent())
		//http.Redirect(w, r, "/", http.StatusFound)
	}
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	fmt.Fprintf(w, "yasdkf")
}
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

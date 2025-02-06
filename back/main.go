package main

import (
	"fmt"
	"log"
	"net/http"
)

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "yourusername"
// 	password = "yourpassword"
// 	dbname   = "yourdbname"
// )

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World!")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

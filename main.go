package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func main() {
	// Load method reads the .env file and loads the set variables into
	// the environment so that they can be accessed through the os.Getenv()
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	fmt.Println("Serving on port:", port)
	http.ListenAndServe(":"+port, mux)
}

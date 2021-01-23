package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	searchQuery := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	fmt.Println("Search Query is:", searchQuery)
	fmt.Println("Page is:", page)
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

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("search", searchHandler)
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	fmt.Println("Serving on port:", port)
	http.ListenAndServe(":"+port, mux)
}

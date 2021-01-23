package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/News-API-Web-app/news"
	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func searchHandler(newsApi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		results, err := newsApi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v", results)
	}
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

	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}

	myClient := &http.Client{Timeout: 10 * time.Second}
	newsApi := news.NewClient(myClient, apiKey, 20)

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("search", searchHandler(newsApi))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	fmt.Println("Serving on port:", port)
	http.ListenAndServe(":"+port, mux)
}

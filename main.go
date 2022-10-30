package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Bookmark struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type StaticFiles struct {
	IndexHTML string
	MainJS    string
	StyleCSS  string
}

func readFileAsStringOrPanic(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func main() {
	port := "8000"
	if len(os.Getenv("PORT")) > 0 {
		if _, err := strconv.Atoi(os.Getenv("PORT")); err != nil {
			log.Fatal("Invalid environment variable: PORT")
		}
		port = os.Getenv("PORT")
	}

	staticFiles := StaticFiles{
		IndexHTML: readFileAsStringOrPanic("./static/index.html"),
		MainJS:    readFileAsStringOrPanic("./static/main.js"),
		StyleCSS:  readFileAsStringOrPanic("./static/style.css"),
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(staticFiles.IndexHTML))
	}).Methods("GET")

	router.HandleFunc("/main.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/javascript; charset=utf-8")
		w.Write([]byte(staticFiles.MainJS))
	}).Methods("GET")

	router.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css; charset=utf-8")
		w.Write([]byte(staticFiles.StyleCSS))
	}).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

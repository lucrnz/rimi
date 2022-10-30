package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
)

type BookmarkStore struct {
	mu   sync.Mutex
	data []Bookmark
}

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

	store := BookmarkStore{
		data: make([]Bookmark, 0),
	}

	store.mu.Lock()
	if _, err := os.Stat("./data.json"); err == nil {
		data, err := os.ReadFile("./data.json")
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(data, &store.data)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Loaded %v bookmarks.\n", len(store.data))
	store.mu.Unlock()

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

	router.HandleFunc("/api/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			msg := fmt.Sprintf("Cannot read request: %v", err)
			fmt.Fprintf(os.Stderr, "%v", msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}

		var bookmark Bookmark
		err = json.Unmarshal(data, &bookmark)
		if err != nil {
			msg := fmt.Sprintf("Cannot parse request: %v", err)
			fmt.Fprintf(os.Stderr, "%v", msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}

		if len(bookmark.Title) > 0 && len(bookmark.URL) > 0 {
			store.mu.Lock()
			defer store.mu.Unlock()
			store.data = append(store.data, bookmark)
			w.WriteHeader(http.StatusCreated)
		} else {
			msg := "Invalid bookmark"
			fmt.Fprintf(os.Stderr, "%v", msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}
	}).Methods("POST")

	router.HandleFunc("/api/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		store.mu.Lock()
		defer store.mu.Unlock()
		data, err := json.Marshal(store.data)
		if err != nil {
			msg := fmt.Sprintf("Cannot serialize store: %v", err)
			fmt.Fprintf(os.Stderr, "%v", msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}
		w.Header().Add("Content-Type", "text/json; charset=UTF-8")
		w.Write(data)
	}).Methods("GET")

	cancelChan := make(chan os.Signal, 1)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		fmt.Printf("Server listening on port %v", port)
		log.Fatal(http.ListenAndServe(":"+port, router))
	}()
	sig := <-cancelChan
	log.Printf("Caught signal %v\n", sig)
	store.mu.Lock()
	defer store.mu.Unlock()
	data, err := json.MarshalIndent(store.data, "", "\t")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("./data.json", data, 0644)
	if err != nil {
		panic(err)
	}
}

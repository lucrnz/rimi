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
	"strings"
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
	IndexHTML  string
	MainJS     string
	StyleCSS   string
	FavIconSVG string
}

func readFileAsStringOrPanic(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func main() {
	bind := "localhost"
	if len(os.Getenv("BIND")) > 0 {
		bind = os.Getenv("BIND")
	}

	port := "8000"
	if len(os.Getenv("PORT")) > 0 {
		if _, err := strconv.Atoi(os.Getenv("PORT")); err != nil {
			log.Fatal("Invalid environment variable: PORT")
		}
		port = os.Getenv("PORT")
	}

	title := "rimi - bookmark manager"
	if len(os.Getenv("TITLE")) > 0 {
		title = os.Getenv("TITLE")
	}

	staticFiles := StaticFiles{
		IndexHTML:  strings.ReplaceAll(readFileAsStringOrPanic("./static/index.html"), "%TITLE%", title),
		MainJS:     readFileAsStringOrPanic("./static/main.js"),
		StyleCSS:   readFileAsStringOrPanic("./static/style.css"),
		FavIconSVG: readFileAsStringOrPanic("./static/favicon.svg"),
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

	router.HandleFunc("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/svg+xml; charset=utf-8")
		w.Write([]byte(staticFiles.FavIconSVG))
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

	router.HandleFunc("/api/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		requestInfo := struct {
			URL string `json:"url"`
		}{}
		err := decoder.Decode(&requestInfo)
		if err != nil {
			msg := fmt.Sprintf("Cannot decode request: %v", err)
			fmt.Fprintf(os.Stderr, "%v", msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}

		if len(requestInfo.URL) > 0 {
			store.mu.Lock()
			defer store.mu.Unlock()

			elIdx := -1
			for idx, el := range store.data {
				if el.URL == requestInfo.URL {
					elIdx = idx
					break
				}
			}
			if elIdx > -1 {
				// remove element from array using slicing
				// read: https://stackoverflow.com/a/57213476
				newData := make([]Bookmark, 0)
				newData = append(newData, store.data[:elIdx]...)
				store.data = append(newData, store.data[elIdx+1:]...)
				w.WriteHeader(http.StatusOK)
			} else {
				msg := "Bookmark not found"
				fmt.Fprintf(os.Stderr, "%v", msg)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(msg))
			}
		} else {
			msg := "Invalid bookmark"
			fmt.Fprintf(os.Stderr, "%v", msg)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(msg))
			return
		}
	}).Methods("DELETE")

	cancelChan := make(chan os.Signal, 1)
	// catch SIGETRM or SIGINTERRUPT
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		fmt.Printf("rimi will listen on %v:%v\n", bind, port)
		log.Fatal(http.ListenAndServe(bind+":"+port, router))
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

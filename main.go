package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// mapping the short url string with the whole URL struct object for that short url
var urlDB = make(map[string]URL)

// generates a hash for the original url string
func generateURLHash(originalUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalUrl))

	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)

	return hash[:8]
}

// creates a short url and stores its data in the DB
func createShortURL(originalUrl string) string {
	shortUrl := generateURLHash(originalUrl)
	id := shortUrl

	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalUrl,
		ShortURL:     shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}

// retrieves original URL from the shortened URL
func getOriginalURL(id string) (URL, error) {
	url, ok := urlDB[id]

	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL := createShortURL(data.URL)

	response := struct {
		ShortenedURL string `json:"short_url"`
	}{ShortenedURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, "Error while encoding json", http.StatusInternalServerError)
		return
	}
}

func RedirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getOriginalURL(id)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func RootPageUrlHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL Shortener!")
}

func main() {
	fmt.Println("Starting URL shortener...")

	http.HandleFunc("/", RootPageUrlHandler)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", RedirectURLHandler)

	fmt.Println("Starting the server on PORT 3000")
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error starting the server", err)
	}
}

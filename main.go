package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"math/rand"
	"strconv"
)

type Book struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000))
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func updateBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var book Book

	for idx, item := range books {
		if item.ID == params["id"] {
			book = books[idx]
			books = append(books[:idx], books[idx+1:]...)

			_ = json.NewDecoder(r.Body).Decode(&book)
			books = append(books, book)

			break
		}
	}

	json.NewEncoder(w).Encode(book)

}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for idx, item := range books {
		if item.ID == params["id"] {
			books = append(books[:idx], books[idx+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(books)
}

func main() {
	// init router
	r := mux.NewRouter()

	// mock data
	books = append(books, Book{ID: "1", Isbn: "111", Title: "Book One", Author: &Author{FirstName: "John", LastName: "Smith"}})
	books = append(books, Book{ID: "2", Isbn: "222", Title: "Book Two", Author: &Author{FirstName: "Lara", LastName: "Croft"}})

	// route handlers / endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBooks).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBooks).Methods("DELETE")

	fmt.Println("Server is running at localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
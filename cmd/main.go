package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"bookmarkSearch/internal/db"
	"bookmarkSearch/internal/elasticsearch"
	"bookmarkSearch/internal/handler"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	host := os.Getenv("POSTGRES_HOST")
	username := os.Getenv("POSTGRES_USER")
	port := os.Getenv("POSTGRES_PORT")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, dbname)

	db.InitDB(dataSourceName)

	elasticsearch.InitES()

	r := mux.NewRouter()
	r.HandleFunc("/bookmarks", handler.SearchBookmarks).Methods("GET")
	r.HandleFunc("/bookmarks/add", handler.AddBookmark).Methods("POST")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"bookmarkSearch/internal/db"
	"bookmarkSearch/internal/elasticsearch"
)

type Bookmark struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type SearchResponse struct {
	PostgresResults []Bookmark             `json:"postgres_results"`
	ESResults       map[string]interface{} `json:"elasticsearch_results"`
}

func SearchBookmarks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	var wg sync.WaitGroup

	response := SearchResponse{}
	var errOccurred bool

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := db.DB.Query("SELECT id, content FROM bookmarks WHERE content ILIKE $1", "%"+query+"%")
		if err != nil {
			errOccurred = true
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bookmarks []Bookmark
		for rows.Next() {
			var b Bookmark
			if err := rows.Scan(&b.ID, &b.Content); err != nil {
				errOccurred = true
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			bookmarks = append(bookmarks, b)
		}
		response.PostgresResults = bookmarks
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		queryBody := `{"query": {"match": {"content": "` + query + `"}}}`
		res, err := elasticsearch.ES.Search(
			elasticsearch.ES.Search.WithContext(context.Background()),
			elasticsearch.ES.Search.WithIndex("bookmarks"),
			elasticsearch.ES.Search.WithBody(strings.NewReader(queryBody)),
			elasticsearch.ES.Search.WithTrackTotalHits(true),
		)
		if err != nil {
			errOccurred = true
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		var esResults map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&esResults); err != nil {
			errOccurred = true
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.ESResults = esResults
	}()

	wg.Wait()

	if errOccurred {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func AddBookmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var bookmark Bookmark
	if err := json.NewDecoder(r.Body).Decode(&bookmark); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	errChan := make(chan error, 2)

	go func() {
		_, err := db.DB.Exec("INSERT INTO bookmarks (content) VALUES ($1)", bookmark.Content)
		if err != nil {
			log.Printf("Database error: %v", err)
			errChan <- err
			return
		}
		errChan <- nil
	}()

	go func() {
		queryBody, err := json.Marshal(map[string]string{"content": bookmark.Content})
		if err != nil {
			log.Printf("Error marshaling JSON for Elasticsearch: %v", err)
			errChan <- err
			return
		}

		res, err := elasticsearch.ES.Index(
			"bookmarks",
			strings.NewReader(string(queryBody)),
			elasticsearch.ES.Index.WithRefresh("true"),
		)
		if err != nil {
			log.Printf("Elasticsearch error: %v", err)
			errChan <- err
			return
		}
		defer res.Body.Close()

		if res.IsError() {
			var esError map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&esError); err != nil {
				log.Printf("Error decoding Elasticsearch response: %v", err)
				errChan <- err
				return
			}
			errChan <- fmt.Errorf("Elasticsearch error: %v", esError)
			return
		}

		errChan <- nil
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bookmark)
}

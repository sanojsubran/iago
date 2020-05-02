package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
1. fetch the storyIds
2. fetch the story details for each story Ids
3. process it to list of links
*/

type server struct{}

/*
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}*/

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

func main() {
	//story_ids_url := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	//story_detail_url := "https://hacker-news.firebaseio.com/v0/item/story_identifier.json?print=pretty"
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", r))
}

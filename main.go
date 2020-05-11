package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

/*
1. fetch the storyIds
2. fetch the story details for each story Ids
3. process it to list of links
*/

// type server struct{}

/*
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}*/
type StoryID struct {
	Array []int64
}

func getHackerNews() {
	storyIdsURL := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	//storyDetailURL := "https://hacker-news.firebaseio.com/v0/item/story_identifier.json?print=pretty"
	storyIDJSON, err := http.Get(storyIdsURL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	topstoris, err := ioutil.ReadAll(storyIDJSON.Body)
	if err != nil {
		log.Fatal(err)
	}
	//storyList := string(topstoris)
	arr := StoryID{}
	err = json.Unmarshal([]byte(topstoris), &arr.Array)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log.Printf("Unmarshaled: %v", arr[0:20])
	//storyID := "23143888"

	//storyDetailURL = strings.Replace(storyDetailURL, "story_identifier", storyID, -1)

}

// func home(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message": "hello world"}`))
// }

func main() {
	//story_ids_url := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	//story_detail_url := "https://hacker-news.firebaseio.com/v0/item/story_identifier.json?print=pretty"
	// r := mux.NewRouter()
	// r.HandleFunc("/", home)
	// log.Fatal(http.ListenAndServe(":8080", r))
	getHackerNews()
}

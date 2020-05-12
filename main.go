package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//StoryID ... exported story
type StoryID struct {
	Array []int64
}

//Story ... exported story
type Story struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func getHackerNews() {
	storyIdsURL := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	storyDetailURL := "https://hacker-news.firebaseio.com/v0/item/story_identifier.json?print=pretty"
	storyIDJSON, err := http.Get(storyIdsURL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	topstoris, err := ioutil.ReadAll(storyIDJSON.Body)
	if err != nil {
		log.Fatal(err)
	}
	arr := StoryID{}
	err = json.Unmarshal([]byte(topstoris), &arr.Array)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	top20StoryIDs := arr.Array[:20]
	fmt.Printf("Top 20 stories: %v\n", top20StoryIDs)
	for _, story := range top20StoryIDs {
		url := strings.Replace(storyDetailURL, "story_identifier", strconv.FormatInt(story, 10), -1)
		//fmt.Printf("Story url: %v\n", url)
		detailsReq, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		details, err := ioutil.ReadAll(detailsReq.Body)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		story := Story{}
		err = json.Unmarshal([]byte(details), &story)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("Story contents: %v | %v | %v\n", story.Id, story.Title, story.Url)
	}
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

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

	"github.com/gorilla/mux"
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

func getHackerNews(w http.ResponseWriter, r *http.Request) {
	storyIdsURL := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	storyDetailURL := "https://hacker-news.firebaseio.com/v0/item/story_identifier.json?print=pretty"
	storyIDJSON, err := http.Get(storyIdsURL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer storyIDJSON.Body.Close()
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
	top20StoryIDs := arr.Array[:30]
	HNData := make([]Story, 0)
	fmt.Printf("Top 20 stories: %v\n", top20StoryIDs)
	for _, story := range top20StoryIDs {
		url := strings.Replace(storyDetailURL, "story_identifier", strconv.FormatInt(story, 10), -1)
		//fmt.Printf("Story url: %v\n", url)
		detailsReq, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer detailsReq.Body.Close()
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
		HNData = append(HNData, story)
		//fmt.Printf("Story contents: %v | %v | %v\n", story.ID, story.Title, story.URL)
		//w.Write([]byte(details))
	}
	// for _, ekStory := range HNData {
	// 	fmt.Println(ekStory)
	// }
	marshalledData, _ := json.Marshal(HNData)
	fmt.Println(string(marshalledData))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshalledData)
}

func main() {
	//getHackerNews()
	r := mux.NewRouter()
	//s := r.Host("www.localhost").Subrouter()

	r.HandleFunc("/", getHackerNews)
	log.Fatal(http.ListenAndServe(":8080", r))
}

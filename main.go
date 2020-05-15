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
	"time"

	"github.com/gorilla/handlers"
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

type HN struct {
	HNJsonData []byte
}

func (a *HN) getHackerNews() {
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
		if len(story.URL) == 0 {
			story.URL = string("https://news.ycombinator.com/item?id=") + strconv.FormatInt(story.ID, 10)
		}
		HNData = append(HNData, story)
		//fmt.Printf("Story contents: %v | %v | %v\n", story.ID, story.Title, story.URL)
		//w.Write([]byte(details))
	}
	// for _, ekStory := range HNData {
	// 	fmt.Println(ekStory)
	// }
	a.HNJsonData, _ = json.Marshal(HNData)
	//return marshalledData
}

//HandleHackerNews...
func (a *HN) HandleHackerNews(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
	fmt.Println("Hello")
	w.Write(a.HNJsonData)
}

func main() {
	//getHackerNews()
	var obj HN

	go func() {
		for {
			obj.getHackerNews()
			time.Sleep(5 * time.Minute)
		}
	}()
	r := mux.NewRouter()
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	//s := r.Host("www.localhost").Subrouter()
	r.HandleFunc("/", obj.HandleHackerNews)
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(header, methods, origins)(r)))
}

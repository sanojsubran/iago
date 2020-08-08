package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type NewsRequest struct {
	mutex       *sync.Mutex
	newsContent map[string][]storyEntry
}

func (n *NewsRequest) generateNewsStream() []byte {
	n.mutex.Lock()
	data, err := json.Marshal(n.newsContent)
	n.mutex.Unlock()
	if nil != err {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
	return data
}

func (n *NewsRequest) handleNewsReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a request. Processing...")
	data := n.generateNewsStream()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (n *NewsRequest) updateNewsFeed(src string, data []storyEntry) {
	n.mutex.Lock()
	n.newsContent[src] = data
	n.mutex.Unlock()
}

func main() {
	r := mux.NewRouter()
	news := NewsRequest{}

	news.mutex = &sync.Mutex{}
	news.newsContent = make(map[string][]storyEntry)

	hn := hackerNews{"hacker_news"}
	rdpgm := redditPgm{"reddit_pgm"}
	rdcpp := redditCpp{"reddit_cpp"}
	godev := golangDev{"golang_dev"}
	reactdev := reactDev{"react_dev"}
	tcrunch := techCrunch{"techcrunch"}
	slashdot := slashDot{"slashdot"}

	go func() {
		for {
			src, data := getFeed(rdpgm, 30)
			news.updateNewsFeed(src, data)
			time.Sleep(15 * time.Minute)
		}
	}()

	go func() {
		for {
			src, data := getFeed(hn, 30)
			news.updateNewsFeed(src, data)
			time.Sleep(15 * time.Minute)
		}
	}()

	go func() {
		for {
			src, data := getFeed(rdcpp, 30)
			news.updateNewsFeed(src, data)
			time.Sleep(15 * time.Minute)
		}
	}()

	go func() {
		for {
			src, data := getFeed(godev, 10)
			news.updateNewsFeed(src, data)
			time.Sleep(15 * time.Minute)
		}
	}()

	go func() {
		for {
			src, data := getFeed(reactdev, 10)
			news.updateNewsFeed(src, data)
			time.Sleep(15 * time.Minute)
		}
	}()

	go func() {
		for {
			src, data := getFeed(tcrunch, 10)
			news.updateNewsFeed(src, data)
			time.Sleep(15 * time.Minute)
		}
	}()

	go func() {
		for {
			src, data := getFeed(slashdot, 10)
			news.updateNewsFeed(src, data)
			time.Sleep(15 * time.Minute)
		}
	}()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.HandleFunc("/", news.handleNewsReq)
	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(header, methods, origins)(r)))
}

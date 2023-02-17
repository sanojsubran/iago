package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
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
		fmt.Println("Unable to marshal the json data. Error: " + err.Error())
		return []byte{}
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

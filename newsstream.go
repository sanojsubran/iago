package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type newsStream struct {
	mutex       *sync.Mutex
	newsContent map[string][]story
}

func (n *newsStream) generateNewsStream() []byte {
	n.mutex.Lock()
	data, err := json.Marshal(n.newsContent)
	n.mutex.Unlock()

	if nil != err {
		logrus.Error("unable to marshal the json data. Error: " + err.Error())
		return nil
	}

	return data
}

func (n *newsStream) handleNewsReq(w http.ResponseWriter, r *http.Request) {
	logrus.Info("received a request. Processing...")

	data := n.generateNewsStream()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)
	if wC, err := w.Write(data); err != nil {
		logrus.Errorf("\nerror while writing response. error: %s, write count: %d", err.Error(), wC)
	}
}

func (n *newsStream) updateNewsFeed(src string, data []story) {
	n.mutex.Lock()
	n.newsContent[src] = data
	n.mutex.Unlock()
}

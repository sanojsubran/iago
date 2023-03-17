package iago

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

type NewsHandler struct {
	mutex       *sync.Mutex
	newsContent map[string][]story
}

// feedContent represents the complete feed data from a single source.
type feedContent struct {
	Title   string
	Article []story
}

// story represent the single entry in the feed. While RSS entries might not have an ID, other sources might provide
// IDs which might be useful in the future.
type story struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

type NewsSource struct {
	Url        string
	ClientType string
	Topic      string
	EntryCount int
}

func (n *NewsHandler) Init() {
	n.mutex = &sync.Mutex{}
	n.newsContent = make(map[string][]story)
}

// generateNewsStream generates the data to be sent via HTTP from the news content map.
func (n *NewsHandler) generateNewsStream() []byte {
	n.mutex.Lock()
	data, err := json.Marshal(n.newsContent)
	n.mutex.Unlock()

	if nil != err {
		logrus.Error("unable to marshal the json data. Error: " + err.Error())
		return nil
	}

	return data
}

func (n *NewsHandler) getSourceHandler(s NewsSource) fetcher {
	switch s.ClientType {
	case "reddit":
		return redditClient{name: s.Topic, url: s.Url, count: s.EntryCount}
	case "rss":
		return rssClient{name: s.Topic, url: s.Url, count: s.EntryCount}
	case "hackernews":
		return hnClient{name: s.Topic, count: s.EntryCount}
	}
	return nil
}

func (n *NewsHandler) UpdateFeed(s NewsSource) error {
	readClient := n.getSourceHandler(s)

	if readClient == nil {
		return errors.New("unable to identify the client reader for the news source")
	}

	d, err := readClient.readData()

	if err != nil {
		logrus.Error("unable to fetch the from the news source. Error: " + err.Error())
		return err
	}

	n.updateNewsFeed(d.Title, d.Article)

	return nil
}

func (n *NewsHandler) HandleNewsReq(w http.ResponseWriter, r *http.Request) {
	logrus.Info("received a request. Processing...")

	data := n.generateNewsStream()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)
	if wC, err := w.Write(data); err != nil {
		logrus.Errorf("\nerror while writing response. error: %s, write count: %d", err.Error(), wC)
	}
}

func (n *NewsHandler) updateNewsFeed(src string, data []story) {
	n.mutex.Lock()
	n.newsContent[src] = data
	n.mutex.Unlock()
}

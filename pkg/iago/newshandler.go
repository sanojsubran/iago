package iago

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
)

// Config defines settings of app. Currently, the embedded file provides details of the news sources.
type Config struct {
	Sources []NewsSource `json:"sources"`
}

// NewsHandler defines the type which stores and handles the feed data.
type NewsHandler struct {
	mutex       *sync.Mutex
	newsContent map[string][]story
}

// NewsSource defines the type which stores the details of a news source.
type NewsSource struct {
	Url        string `json:"url"`
	ClientType string `json:"clientType"`
	Topic      string `json:"topic"`
	EntryCount int    `json:"entryCount"`
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

// Init initializes the newsHandler struct.
func (n *NewsHandler) Init() {
	n.mutex = &sync.Mutex{}
	n.newsContent = make(map[string][]story)
}

// GetConfiguredSources gets the news sources from the config data.
func (n *NewsHandler) GetConfiguredSources(config string) []NewsSource {

	var configData Config

	err := json.Unmarshal([]byte(config), &configData)

	if err != nil {
		logrus.Error("unable to parse the config data. Error: ", err.Error())
	}

	sources := configData.Sources

	if len(sources) == 0 {
		logrus.Error("count of parsed sources is zero. Check the config.json.")
	}

	return sources
}

// UpdateFeed updates the news feed data content.
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

// HandleNewsRequests handles the incoming requests.
func (n *NewsHandler) HandleNewsRequests(w http.ResponseWriter, _ *http.Request) {
	logrus.Info("received a request. Processing...")

	data := n.generateNewsStream()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)
	if wC, err := w.Write(data); err != nil {
		logrus.Errorf("\nerror while writing response. error: %s, write count: %d", err.Error(), wC)
	}
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

func (n *NewsHandler) updateNewsFeed(src string, data []story) {
	n.mutex.Lock()
	n.newsContent[src] = data
	n.mutex.Unlock()
}

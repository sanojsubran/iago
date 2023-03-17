package iago

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

// rssFeedType is the data received from the actual feed in json format.
type rssFeedType struct {
	Items []struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	} `json:"items"`
}

// rssClient defines the readData interface method and fetches the data from the rss data source.
type rssClient struct {
	name  string
	url   string
	count int
}

func (rss rssClient) clientName() string {
	return rss.name
}

func (rss rssClient) readData() (feedContent, error) {
	fp := gofeed.NewParser()
	data := feedContent{}

	feed, err := fp.ParseURL(rss.url)
	if err != nil {
		logrus.Warn("error occurred while fetching the feed")
		return data, err
	}

	jsonContent := rssFeedType{}
	err = json.Unmarshal([]byte(feed.String()), &jsonContent)
	if nil != err {
		fmt.Println("Unable to unmarshal the data from TC. Error: ", err.Error())
		return data, err
	}

	stories := make([]story, 0)
	for i, v := range jsonContent.Items {
		stories = append(stories, story{
			Id:    int64(i),
			Title: v.Title,
			Url:   v.Link,
		})
	}

	data = feedContent{
		Title:   rss.name,
		Article: stories[:rss.count],
	}

	return data, nil
}

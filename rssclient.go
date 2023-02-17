package main

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
	url   string
	count int
}

func (c rssClient) readData() (feedContent, error) {
	fp := gofeed.NewParser()
	data := feedContent{}

	feed, err := fp.ParseURL(c.url)
	if err != nil {
		logrus.Warn("error occurred while fetching the feed")
		return data, err
	}

	jsonContent := rssFeedType{}
	err = json.Unmarshal([]byte(feed.String()), &jsonContent)
	if nil != err {
		fmt.Println("Unable to unmars hal the data from TC. Error: ", err.Error())
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
		Title:   c.url,
		Article: stories[:c.count],
	}

	return data, nil
}

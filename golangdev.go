package main

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
)

type golangDev struct {
	newsSrc string
}

func (gd golangDev) readData(count int16) (string, []storyEntry) {
	fp := gofeed.NewParser()
	stories := make([]storyEntry, 0)
	feed, _ := fp.ParseURL("https://blog.golang.org/feed.atom")
	if feed.String() == "" {
		fmt.Println("Empty data received from the source")
		return gd.newsSrc, stories
	}
	jsonContent := RSSFeedContent{}
	err := json.Unmarshal([]byte(feed.String()), &jsonContent)
	if nil != err {
		fmt.Println("Unable to unmarshall the data from golangdev. Error: " + err.Error())
		return gd.newsSrc, stories
	}
	for index, story := range jsonContent.Items {
		stories = append(stories, storyEntry{
			Id:    int64(index),
			Title: story.Title,
			Url:   story.Link,
		})
	}
	return gd.newsSrc, stories[:count]
}

package main

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
)

type reactDev struct {
	newsSrc string
}

func (rd reactDev) readData(count int16) (string, []storyEntry) {
	stories := make([]storyEntry, 0)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://reactjs.org/feed.xml")
	if feed.String() == "" {
		fmt.Println("Empty data received from the source")
		return rd.newsSrc, stories
	}
	jsonContent := RSSFeedContent{}
	err := json.Unmarshal([]byte(feed.String()), &jsonContent)
	if nil != err {
		fmt.Println("Unable to unmarshall the data from reactdev. Error: " + err.Error())
		return rd.newsSrc, stories
	}
	for index, story := range jsonContent.Items {
		stories = append(stories, storyEntry{
			Id:    int64(index),
			Title: story.Title,
			Url:   story.Link,
		})
	}
	return rd.newsSrc, stories[:count]
}

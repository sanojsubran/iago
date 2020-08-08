package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
)

type reactDev struct {
	newsSrc string
}

func (rd reactDev) readData(count int16) (string, []storyEntry) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://reactjs.org/feed.xml")
	jsonContent := RSSFeedContent{}
	err := json.Unmarshal([]byte(feed.String()), &jsonContent)
	if nil != err {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
	stories := make([]storyEntry, 0)
	for index, story := range jsonContent.Items {
		stories = append(stories, storyEntry{
			Id:    int64(index),
			Title: story.Title,
			Url:   story.Link,
		})
	}
	return rd.newsSrc, stories[:count]
}

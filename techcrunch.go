package main

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
)

type techCrunch struct {
	newsSrc string
}

func (tc techCrunch) readData(count int16) (string, []storyEntry) {
	fp := gofeed.NewParser()
	stories := make([]storyEntry, 0)
	feed, _ := fp.ParseURL("https://techcrunch.com/feed/")
	jsonContent := RSSFeedContent{}
	err := json.Unmarshal([]byte(feed.String()), &jsonContent)
	if nil != err {
		fmt.Println("Unable to unmarshal the data from TC. Error: ", err.Error())
		return tc.newsSrc, stories
	}
	for index, story := range jsonContent.Items {
		stories = append(stories, storyEntry{
			Id:    int64(index),
			Title: story.Title,
			Url:   story.Link,
		})
	}
	return tc.newsSrc, stories[:count]
}

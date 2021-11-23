package main

import (
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
)

type slashDot struct {
	newsSrc string
}

func (sd slashDot) readData(count int16) (string, []storyEntry) {
	stories := make([]storyEntry, 0)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://rss.slashdot.org/Slashdot/slashdotMain")
	jsonContent := RSSFeedContent{}
	if feed.String() == "" {
		fmt.Println("Empty data received from the source")
		return sd.newsSrc, stories
	}
	err := json.Unmarshal([]byte(feed.String()), &jsonContent)
	if nil != err {
		fmt.Println("Unable to unmarshal the data from Slashdot. Error: ", err.Error())
		return sd.newsSrc, stories
	}
	for index, story := range jsonContent.Items {
		stories = append(stories, storyEntry{
			Id:    int64(index),
			Title: story.Title,
			Url:   story.Link,
		})
	}
	return sd.newsSrc, stories[:count]
}

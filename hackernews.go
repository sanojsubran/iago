package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type hackerNews struct {
	newsSrc string
}

func (hn hackerNews) readData(count int16) (string, []storyEntry) {
	storyListApi := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	storyDetailsApi := "https://hacker-news.firebaseio.com/v0/item/story_identifier.json?print=pretty"
	storyListJson, err := http.Get(storyListApi)
	if nil != err {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer storyListJson.Body.Close()
	topstories, err := ioutil.ReadAll(storyListJson.Body)
	if nil != err {
		log.Fatal(err)
	}
	storyList := make([]int64, 0)
	err = json.Unmarshal([]byte(topstories), &storyList)
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	topStoryIds := storyList[:count]
	finalStoryList := make([]storyEntry, 0)
	for _, story := range topStoryIds {
		url := strings.Replace(storyDetailsApi, "story_identifier", strconv.FormatInt(story, 10), -1)
		details, err := http.Get(url)
		if nil != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer details.Body.Close()
		detailContent, err := ioutil.ReadAll(details.Body)
		if nil != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		entry := storyEntry{}
		err = json.Unmarshal([]byte(detailContent), &entry)
		if nil != err {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if len(entry.Url) == 0 {
			entry.Url = string("https://news.ycombinator.com/item?id=") + strconv.FormatInt(entry.Id, 10)
		}
		finalStoryList = append(finalStoryList, entry)
	}
	return hn.newsSrc, finalStoryList[:count]
}

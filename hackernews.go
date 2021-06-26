package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	finalStoryList := make([]storyEntry, 0)
	if nil != err {
		fmt.Println("Unable to reach the Url: " + storyDetailsApi + " Error: " + err.Error())
		return hn.newsSrc, finalStoryList
	}
	defer storyListJson.Body.Close()
	topstories, err := ioutil.ReadAll(storyListJson.Body)
	if nil != err {
		fmt.Println("Unable to read the request body from HN. Error: " + err.Error())
		return hn.newsSrc, finalStoryList
	}
	storyList := make([]int64, 0)
	err = json.Unmarshal([]byte(topstories), &storyList)
	if nil != err {
		fmt.Println("Unable to unmarshall the data from HN. Error: " + err.Error())
		return hn.newsSrc, finalStoryList
	}
	topStoryIds := storyList[:count]
	for _, story := range topStoryIds {
		url := strings.Replace(storyDetailsApi, "story_identifier", strconv.FormatInt(story, 10), -1)
		details, err := http.Get(url)
		if nil != err {
			fmt.Println("Unable to fetch the story details from HN. Error: " + err.Error())
			return hn.newsSrc, finalStoryList
		}
		defer details.Body.Close()
		detailContent, err := ioutil.ReadAll(details.Body)
		if nil != err {
			fmt.Println("Unable to read the story details body from HN. Error: ", err.Error())
			return hn.newsSrc, finalStoryList
		}
		entry := storyEntry{}
		err = json.Unmarshal([]byte(detailContent), &entry)
		if nil != err {
			fmt.Println("Unable to unmarshall the story details from Hn. Error: " + err.Error())
			return hn.newsSrc, finalStoryList
		}
		if len(entry.Url) == 0 {
			entry.Url = string("https://news.ycombinator.com/item?id=") + strconv.FormatInt(entry.Id, 10)
		}
		finalStoryList = append(finalStoryList, entry)
	}
	return hn.newsSrc, finalStoryList[:count]
}

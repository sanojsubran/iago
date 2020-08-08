package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type redditCpp struct {
	newsSrc string
}

//JSONContent
type JSONRCPPContent struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Data struct {
				Title string `json:"title"`
				Id    string `json:"id"`
				Url   string `json:"url"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func (rdp redditCpp) readData(count int16) (string, []storyEntry) {
	url := "https://www.reddit.com/r/cpp/new.json?limit=count"
	url = strings.Replace(url, "count", strconv.FormatInt(int64(count), 10), -1)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", "webapp iago")
	resp, err := client.Do(req)
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	story := JSONRCPPContent{}
	err = json.Unmarshal(b, &story)
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	children := story.Data.Children
	stories := make([]storyEntry, 0)
	for index, child := range children {
		stories = append(stories, storyEntry{
			Id:    int64(index),
			Title: child.Data.Title,
			Url:   child.Data.Url})
	}
	return rdp.newsSrc, stories[:count]
}

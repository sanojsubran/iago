package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// hnClient defines the readData interface method and fetches the data from hacker news.
type hnClient struct {
	name  string
	count int
}

func (hn hnClient) clientName() string {
	return hn.name
}

func (hn hnClient) readData() (feedContent, error) {
	storyList, err := hn.getStoryIDList()
	if err != nil {
		logrus.Error("error occurred while fetching the story list from HN")
		return feedContent{}, err
	}

	topStoryIds := storyList[:hn.count]

	finalStoryList := make([]story, 0)

	for _, v := range topStoryIds {
		s, err := hn.getStoryDetail(v)
		if err != nil {
			logrus.Error("unable to fetch the story details. Skipping to next item.")
			continue
		}

		finalStoryList = append(finalStoryList, s)
	}

	data := feedContent{
		Title:   hn.name,
		Article: finalStoryList[:],
	}

	return data, nil
}

func (hn hnClient) getStoryIDList() ([]int64, error) {
	storyListApi := "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"

	storyListJson, err := http.Get(storyListApi)
	defer storyListJson.Body.Close()
	if nil != err {
		fmt.Println("unable to reach the Url: " + storyListApi + " Error: " + err.Error())
		return nil, err
	}

	topStories, err := io.ReadAll(storyListJson.Body)
	if nil != err {
		fmt.Println("unable to read the request body from HN. Error: " + err.Error())
		return nil, err
	}

	storyList := make([]int64, 0)
	err = json.Unmarshal(topStories, &storyList)
	if nil != err {
		fmt.Println("unable to unmarshall the data from HN. Error: " + err.Error())
		return nil, err
	}

	return storyList, nil
}

func (hn hnClient) getStoryDetail(id int64) (story, error) {
	storyDetailsApi := "https://hacker-news.firebaseio.com/v0/item/story_identifier.json?print=pretty"

	url := strings.Replace(storyDetailsApi, "story_identifier", strconv.FormatInt(id, 10), -1)
	details, err := http.Get(url)
	if nil != err {
		fmt.Println("Unable to fetch the story details from HN. Error: " + err.Error())
		return story{}, err
	}

	detailContent, err := io.ReadAll(details.Body)
	if nil != err {
		fmt.Println("Unable to read the story details body from HN. Error: ", err.Error())
		return story{}, err
	}
	details.Body.Close()

	entry := story{}
	err = json.Unmarshal(detailContent, &entry)
	if nil != err {
		fmt.Println("Unable to unmarshall the story details from Hn. Error: " + err.Error())
		return story{}, err
	}

	// Case where the url is a discussion thread in HN
	if len(entry.Url) == 0 {
		entry.Url = string("https://news.ycombinator.com/item?id=") + strconv.FormatInt(entry.Id, 10)
	}

	return entry, nil
}

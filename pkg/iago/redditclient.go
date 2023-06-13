package iago

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

// RedditJsonType represents the json content received from the reddit website.
type RedditJsonType struct {
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

// redditClient defines the readData interface method and fetches the data from reddit.
type redditClient struct {
	name  string
	url   string
	count int
}

func (rd redditClient) clientName() string {
	return rd.name
}

// readData implements the fetching of data from reddit's sub topics
func (rd redditClient) readData() (feedContent, error) {

	url := rd.url + strconv.FormatInt(int64(rd.count), 10)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", "webapp iago")

	resp, err := client.Do(req)
	if nil != err {
		logrus.Error("unable to query the reddit site. Error: " + err.Error())
		return feedContent{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return feedContent{}, errors.New("query to reddit returned HTTP error code:" + strconv.Itoa(resp.StatusCode))
	}

	b, err := io.ReadAll(resp.Body)
	if nil != err {
		logrus.Error("unable to read the body of reddit response. Error: " + err.Error())
		return feedContent{}, err
	}

	jsonContent := RedditJsonType{}
	err = json.Unmarshal(b, &jsonContent)
	if nil != err {
		logrus.Error("unable to unmarshall the data from reddit. Error: " + err.Error())
		return feedContent{}, err
	}

	children := jsonContent.Data.Children

	stories := make([]story, 0)
	for index, child := range children {
		stories = append(stories, story{
			Id:    int64(index),
			Title: child.Data.Title,
			Url:   child.Data.Url})
	}

	data := feedContent{
		Title:   rd.name,
		Article: stories,
	}

	return data, nil
}

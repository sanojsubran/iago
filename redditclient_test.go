package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedditFeed(t *testing.T) {
	testcases := map[string]struct {
		url   string
		count int
	}{
		"reddit_programming": {
			url:   "https://www.reddit.com/r/programming/top.json?limit=",
			count: 10,
		},
		"reddit_cpp": {
			url:   "https://www.reddit.com/r/cpp/top.json?limit=",
			count: 10,
		},
		"reddit_soccer": {
			url:   "https://www.reddit.com/r/soccer/top.json?limit=",
			count: 10,
		},
		"reddit_japan": {
			url:   "https://www.reddit.com/r/japan/top.json?limit=",
			count: 10,
		},
	}

	for name, test := range testcases {
		t.Run(name, func(t *testing.T) {

			client := redditClient{
				url:   test.url,
				count: test.count,
			}

			feed, err := client.readData()

			assert.NoError(t, err)
			assert.NotEmpty(t, feed.Title)
			assert.NotEqualf(t, 0, len(feed.Article), "number of fetched articles cannot be zero")
		})
	}

}

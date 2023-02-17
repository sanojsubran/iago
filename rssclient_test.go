package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRSSFeed(t *testing.T) {
	testcases := map[string]struct {
		name  string
		url   string
		count int
	}{
		"techcrunch": {
			name:  "techcrunch",
			url:   "https://techcrunch.com/feed/",
			count: 10,
		},
		"slashdot": {
			name:  "slashdot",
			url:   "http://rss.slashdot.org/Slashdot/slashdotMain",
			count: 10,
		},
		"reactdev": {
			name:  "reactdev",
			url:   "https://reactjs.org/feed.xml",
			count: 10,
		},
		"godev": {
			name:  "go.dev",
			url:   "https://go.dev/blog/feed.atom",
			count: 10,
		},
	}

	for name, test := range testcases {
		t.Run(name, func(t *testing.T) {

			client := rssClient{
				name:  test.name,
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

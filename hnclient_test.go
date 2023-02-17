package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHNFeed(t *testing.T) {
	testcases := map[string]struct {
		url   string
		count int
	}{
		"hackernews_30": {
			count: 30,
		},
		"hackernews_40": {
			count: 40,
		},
		"hackernews_50": {
			count: 50,
		},
	}

	for name, test := range testcases {
		t.Run(name, func(t *testing.T) {

			client := hnClient{
				count: test.count,
			}

			feed, err := client.readData()

			assert.NoError(t, err)
			assert.NotEqualf(t, 0, len(feed.Article), "number of fetched articles cannot be zero")
			assert.Equal(t, test.count, len(feed.Article))
		})
	}

}

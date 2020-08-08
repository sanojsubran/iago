package main

type storyEntry struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

type RSSFeedContent struct {
	Items []struct {
		Title string `json:"title"`
		Link  string `json:"link"`
	} `json:"items"`
}

type fetcher interface {
	readData(count int16) (string, []storyEntry)
}

func getFeed(f fetcher, count int16) (string, []storyEntry) {
	src, data := f.readData(count)
	//fmt.Printf("src: %v, data: %+v", src, data)
	return src, data
}

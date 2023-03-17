package iago

type fetcher interface {
	readData() (feedContent, error)
	clientName() string
}

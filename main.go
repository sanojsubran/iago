package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	iago "github.com/sanojsubran/iago/pkg/iago"
)

//go:embed config.json
var cfg string

func main() {
	r := mux.NewRouter()
	var news iago.NewsHandler

	news.Init()

	sources := news.GetConfiguredSources(cfg)

	go func() {
		for {
			for _, source := range sources {
				err := news.UpdateFeed(source)
				if err != nil {
					logrus.Errorf("unable to fetch the data from %s at %d", source.Topic, time.Now().Unix())
				}
			}
			time.Sleep(60 * time.Minute)
		}
	}()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.HandleFunc("/", news.HandleNewsRequests)
	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(header, methods, origins)(r)))

}

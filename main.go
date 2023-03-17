package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/sirupsen/logrus"

	iago "github.com/sanojsubran/iago/pkg/iago"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func readConfig() {
	// Load JSON config.s
	if err := k.Load(file.Provider("./config.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
}

func main() {
	//TODO To be implemented later
	//readConfig()

	sources := make([]iago.NewsSource, 0)

	//TODO: populate the sources with configured news sources

	r := mux.NewRouter()
	var news iago.NewsHandler

	news.Init()

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
	r.HandleFunc("/", news.HandleNewsReq)
	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(header, methods, origins)(r)))

}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/hawx/mpd-scrobbler/client"
	"github.com/hawx/mpd-scrobbler/scrobble"
)

const (
	// only submit tracks longer then minTrackLen
	minTrackLen = 30

	// only submit if played for submitTime second or submitPercentage of length
	submitTime       = 240
	submitPercentage = 50

	// polling interval
	sleepTime = 5 * time.Second
)

const helpMessage = `Usage: mpd-scrobbler [options]

  Scrobbles tracks from mpd.

    --config <path>  # Path to config file (default: './config.toml')
    --db <path>      # Path to database for caching (default: './scrobble.db')
    --port <port>    # Port mpd running on (default: '6600')
    --help           # Display this message
`

var (
	config = flag.String("config", "./config.toml", "")
	dbPath = flag.String("db", "./scrobble.db", "")
	port   = flag.String("port", "6600", "")
	help   = flag.Bool("help", false, "")
)

func catchInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	s := <-c
	log.Printf("caught %s: shutting down", s)
}

func main() {
	flag.Parse()

	if *help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	c, err := client.Dial("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	db, err := scrobble.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile(*config, &conf); err != nil {
		log.Fatal(err)
	}

	apis := []scrobble.Scrobbler{}
	for k, v := range conf {
		api, err := scrobble.New(db, k, v["key"], v["secret"], v["username"], v["password"], v["uri"])
		if err != nil {
			log.Fatal(k, " ", err)
		}

		apis = append(apis, api)
	}

	toSubmit := make(chan client.Song)
	nowPlaying := make(chan client.Song)

	go c.Watch(sleepTime, toSubmit, nowPlaying)
	go func() {
		for {
			select {
			case s := <-nowPlaying:
				for _, api := range apis {
					err := api.NowPlaying(s.Artist, s.Album, s.AlbumArtist, s.Title)
					if err != nil {
						log.Printf("[%s] err(NowPlaying): %s\n", api.Name(), err)
					}
				}

			case s := <-toSubmit:
				for _, api := range apis {
					err := api.Scrobble(s.Artist, s.Album, s.AlbumArtist, s.Title, s.Start)
					if err != nil {
						log.Printf("[%s] err(Scrobble): %s\n", api.Name(), err)
					}
				}
			}
		}
	}()

	catchInterrupt()
}

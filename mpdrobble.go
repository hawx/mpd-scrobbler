package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/hawx/mpdrobble/client"
	"github.com/hawx/mpdrobble/scrobble"
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

const helpMessage = `Usage: mpdrobble [options]

  Scrobbles tracks from mpd.

    --config <path>  # Path to config file (default: './config.toml')
    --port <port>    # Port mpd running on (default: '6600')
    --help           # Display this message
`

var (
	config = flag.String("config", "./config.toml", "")
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

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile(*config, &conf); err != nil {
		log.Fatal(err)
	}

	apis := []*scrobble.Api{}
	for k, v := range conf {
		api, err := scrobble.New(k, v["key"], v["secret"], v["username"], v["password"], v["uri"])
		if err != nil {
			log.Fatal(k, " ", err)
		}

		apis = append(apis, api)
		log.Println("connected to", k)
	}

	toSubmit := make(chan client.PosSong)

	go c.Watch(sleepTime, toSubmit)
	go func() {
		for s := range toSubmit {
			for _, api := range apis {
				err := api.Scrobble(s.Artist, s.Album, s.AlbumArtist, s.Title, s.Start)
				if err != nil {
					log.Println("err(Submit):", err)
				}
			}
			log.Println("Submitted:", s)
		}
	}()

	catchInterrupt()
}

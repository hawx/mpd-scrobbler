package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/hawx/mpdrobble/gompd/mpd"
	"github.com/hawx/mpdrobble/lastfm-go/lastfm"
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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func submitTrack(api *Api, song Song, pos Pos) {
	start := time.Now().UTC().Add(-pos.Duration())

	errs := api.Scrobble(song.artist, song.album, song.albumArtist, song.title, start)
	for err := range errs {
		if err != nil {
			log.Println("submit:", err)
		}
	}
	log.Println("submitted", song)
}

type Pos struct {
	percent float64
	seconds int
	length  int
}

func (pos Pos) Duration() time.Duration {
	return time.Duration(pos.seconds) * time.Second
}

type Song struct {
	title, artist, album, albumArtist, file string
}

func (s Song) String() string {
	return fmt.Sprintf("%s by %s", s.title, s.artist)
}

type Client struct {
	client *mpd.Client
}

func NewClient(client *mpd.Client, err error) (*Client, error) {
	if err != nil {
		return nil, err
	}

	return &Client{client}, err
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Current() (song Song, pos Pos, err error) {
	s, err := c.client.CurrentSong()
	if err != nil {
		return
	}

	st, err := c.client.Status()
	if err != nil {
		return
	}

	parts := strings.Split(st["time"], ":")

	pos.seconds, err = strconv.Atoi(parts[0])
	if err != nil {
		return
	}
	pos.length, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	pos.percent = float64(pos.seconds) / float64(pos.length) * 100

	song = Song{s["Title"], s["Artist"], s["Album"], s["AlbumArtist"], s["file"]}
	return
}

type Api struct {
	apis map[string]*lastfm.Api
}

func NewApi() *Api {
	return &Api{map[string]*lastfm.Api{}}
}

func (apis *Api) Append(name, apiKey, secret, username, password, uriBase string) error {
	api := lastfm.New(apiKey, secret, uriBase)
	err := api.Login(username, password)
	if err != nil {
		return err
	}

	apis.apis[name] = api
	return nil
}

type Err struct {
	name string
	err  error
}

func (e *Err) Error() string {
	return fmt.Sprintf("%v: %v", e.name, e.err.Error())
}

func (apis *Api) Scrobble(artist, album, albumArtist, title string, timestamp time.Time) <-chan error {
	errc := make(chan error, 1)

	go func(artist, album, albumArtist, title string, timestamp time.Time) {
		for name, api := range apis.apis {
			p := lastfm.P{
				"artist":      artist,
				"album":       album,
				"albumArtist": albumArtist,
				"track":       title,
				"timestamp":   timestamp.Unix(),
			}

			_, err := api.Track.Scrobble(p)
			if err != nil {
				errc <- &Err{name, err}
			}
		}
		close(errc)
	}(artist, album, albumArtist, title, timestamp)

	return errc
}

func main() {
	flag.Parse()

	if *help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	c, err := NewClient(mpd.Dial("tcp", ":"+*port))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	api := NewApi()

	var conf map[string]map[string]string
	if _, err := toml.DecodeFile(*config, &conf); err != nil {
		log.Fatal(err)
	}

	for k, v := range conf {
		err = api.Append(k, v["key"], v["secret"], v["username"], v["password"], v["uri"])
		if err != nil {
			log.Fatal(k, " ", err)
		}
		log.Println("connected to", k)
	}

	go func() {
		lastSong := Song{file: ""}
		lastPos := Pos{}
		submitThis := false

		for {
			<-time.After(sleepTime)

			song, pos, err := c.Current()
			if err != nil {
				log.Println(err)
				continue
			}

			if lastSong.file == "" || // have we just started playing?
				song.file != lastSong.file || // are we playing a different track?
				(pos.seconds < lastPos.seconds && pos.Duration() <= sleepTime) { // have we restarted this track?

				log.Println("detected", song)
				if song.title != "" && pos.length >= minTrackLen {
					submitThis = true
					lastSong = song
				}
			} else {
				if submitThis {
					// allow 1 second for rounding errors
					if abs(int(pos.Duration()-lastPos.Duration()-sleepTime)) > int(time.Second) {
						log.Println("Seeking detected, will not submit this track")
						submitThis = false
					} else {
						if pos.seconds >= submitTime || pos.percent >= submitPercentage {
							log.Println("submitting", song)
							go submitTrack(api, song, pos)
							submitThis = false
						}
					}
				}
			}
			lastPos = pos
		}
	}()

	catchInterrupt()
}

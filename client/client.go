package client

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hawx/mpdrobble/client/gompd/mpd"
)

const (
	// only submit if played for submitTime second or submitPercentage of length
	submitTime       = 240
	submitPercentage = 50
)

type PosSong struct {
	Artist      string
	Album       string
	AlbumArtist string
	Title       string
	Start       time.Time
}

type Client struct {
	client    *mpd.Client
	song      mpdSong
	pos       mpdPos
	start     int // stats curtime
	starttime time.Time
	submitted bool
}

func Dial(network, addr string) (*Client, error) {
	c, err := mpd.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	return &Client{c, mpdSong{}, mpdPos{}, 0, time.Now(), false}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Watch(interval time.Duration, ch chan PosSong) {
	for _ = range time.Tick(interval) {
		playtime, err := c.playTime()
		if err != nil {
			log.Println("err(PlayTime):", err)
			continue
		}

		song, err := c.currentSong()
		if err != nil {
			log.Println("err(CurrentSong):", err)
			continue
		}

		pos, err := c.currentPos()
		if err != nil {
			log.Println("err(CurrentPos):", err)
			continue
		}

		// new song
		if song != c.song {
			c.song = song
			c.pos = pos
			c.start = playtime
			c.starttime = time.Now().UTC()
			c.submitted = false
			log.Println("New Song:", song)
		}

		// still playing
		if pos != c.pos {
			c.pos = pos
			if c.canSubmit(playtime) {
				log.Println("Submitting:", song)
				c.submitted = true
				ch <- PosSong{
					Album:       c.song.Album,
					Artist:      c.song.Artist,
					AlbumArtist: c.song.AlbumArtist,
					Title:       c.song.Title,
					Start:       c.starttime,
				}
			}
		}
	}
}

func (c *Client) canSubmit(playtime int) bool {
	if c.submitted || c.song.Artist == "" || c.song.Title == "" {
		return false
	}

	return playtime-c.start >= submitTime ||
		playtime-c.start >= c.pos.Length/(100/submitPercentage)
}

func (c *Client) playTime() (int, error) {
	s, err := c.client.Stats()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(s["playtime"])
}

func (c *Client) currentSong() (mpdSong, error) {
	s, err := c.client.CurrentSong()
	if err != nil {
		return mpdSong{}, nil
	}

	return mpdSong{s["Title"], s["Artist"], s["Album"], s["AlbumArtist"], s["file"]}, nil
}

func (c *Client) currentPos() (pos mpdPos, err error) {
	st, err := c.client.Status()
	if err != nil {
		return
	}

	parts := strings.Split(st["time"], ":")

	pos.Seconds, err = strconv.Atoi(parts[0])
	if err != nil {
		return
	}
	pos.Length, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	pos.Percent = float64(pos.Seconds) / float64(pos.Length) * 100
	return
}

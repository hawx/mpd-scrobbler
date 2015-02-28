package client

import "time"

type Song struct {
	Artist      string
	Album       string
	AlbumArtist string
	Title       string
	Start       time.Time
}

func (s Song) String() string {
	return s.Title + " by " + s.Artist
}

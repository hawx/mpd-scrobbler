package scrobble

import (
	"time"

	"github.com/hawx/mpdrobble/scrobble/lastfm-go/lastfm"
)

type Api struct {
	api  *lastfm.Api
	Name string
}

func New(name, apiKey, secret, username, password, uriBase string) (*Api, error) {
	api := lastfm.New(apiKey, secret, uriBase)
	err := api.Login(username, password)
	if err != nil {
		return nil, err
	}

	return &Api{api, name}, nil
}

func (api *Api) Scrobble(artist, album, albumArtist, title string, timestamp time.Time) error {
	p := lastfm.P{
		"artist":      artist,
		"album":       album,
		"albumArtist": albumArtist,
		"track":       title,
		"timestamp":   timestamp.Unix(),
	}

	_, err := api.api.Track.Scrobble(p)
	return err
}

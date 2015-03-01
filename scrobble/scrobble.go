package scrobble

import (
	"time"

	"github.com/hawx/mpdrobble/scrobble/lastfm"
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
	err := api.api.Scrobble(lastfm.ScrobbleArgs{
		Artist:      artist,
		Album:       album,
		AlbumArtist: albumArtist,
		Track:       title,
		Timestamp:   timestamp.Unix(),
	})
	return err
}

func (api *Api) NowPlaying(artist, album, albumArtist, title string) error {
	err := api.api.UpdateNowPlaying(lastfm.UpdateNowPlayingArgs{
		Artist:      artist,
		Track:       title,
		Album:       album,
		AlbumArtist: albumArtist,
	})
	return err
}

package scrobble

import (
	"log"
	"time"

	"hawx.me/code/mpd-scrobbler/scrobble/lastfm"
)

type Scrobbler interface {
	Scrobble(artist, album, albumArtist, title string, timestamp time.Time) error
	NowPlaying(artist, album, albumArtist, title string) error
	Name() string
}

func New(db Database, name, apiKey, secret, username, password, uriBase string) (Scrobbler, error) {
	api := lastfm.New(apiKey, secret, uriBase)

	queue, err := db.Queue([]byte(name))
	if err != nil {
		return nil, err
	}

	scrobbler := &lastfmScrobbler{api, username, password, false, name}

	log.Printf("[%s] Emptying queue\n", name)
	for {
		track, err := queue.Dequeue()
		if err != nil {
			if err != QUEUE_EMPTY {
				queue.Enqueue(track)
				log.Printf("[%s] Queued: %s by %s\n", name, track.Title, track.Artist)
			}
			log.Printf("[%s] %s\n", name, err)
			break
		}

		err = scrobbler.Scrobble(track.Artist, track.Album, track.AlbumArtist, track.Title, track.Timestamp)
		if err != nil {
			queue.Enqueue(track)
			log.Printf("[%s] Queued: %s by %s\n", name, track.Title, track.Artist)
			log.Printf("[%s] %s\n", name, err)
			break
		}
	}

	log.Printf("[%s] Emptying done\n", name)

	return &queuedScrobbler{scrobbler, queue}, nil
}

type queuedScrobbler struct {
	Scrobbler
	queue Queue
}

func (api *queuedScrobbler) Scrobble(artist, album, albumArtist, title string, timestamp time.Time) (err error) {
	track, err := Track{artist, album, albumArtist, title, timestamp}, nil
	for err == nil {
		err = api.Scrobbler.Scrobble(track.Artist, track.Album, track.AlbumArtist, track.Title, track.Timestamp)
		if err != nil {
			break
		}
		track, err = api.queue.Dequeue()
	}

	if err != nil {
		if err == QUEUE_EMPTY {
			return nil
		} else {
			api.queue.Enqueue(track)
			log.Printf("[%s] Queued: %s by %s\n", api.Name(), title, artist)
		}
	}

	return err
}

type lastfmScrobbler struct {
	api      *lastfm.Api
	username string
	password string
	loggedIn bool
	name     string
}

func (api *lastfmScrobbler) Name() string {
	return api.name
}

func (api *lastfmScrobbler) login() error {
	if !api.loggedIn {
		err := api.api.Login(api.username, api.password)
		if err == nil {
			log.Printf("[%s] Connected", api.Name())
			api.loggedIn = true
		}
		return err
	}
	return nil
}

func (api *lastfmScrobbler) Scrobble(artist, album, albumArtist, title string, timestamp time.Time) error {
	if err := api.login(); err != nil {
		return err
	}

	err := api.api.Scrobble(lastfm.ScrobbleArgs{
		Artist:      artist,
		Album:       album,
		AlbumArtist: albumArtist,
		Track:       title,
		Timestamp:   timestamp.Unix(),
	})

	if err == nil {
		log.Printf("[%s] Submitted: %s by %s\n", api.Name(), title, artist)
	}

	return err
}

func (api *lastfmScrobbler) NowPlaying(artist, album, albumArtist, title string) error {
	if err := api.login(); err != nil {
		return err
	}

	err := api.api.UpdateNowPlaying(lastfm.UpdateNowPlayingArgs{
		Artist:      artist,
		Track:       title,
		Album:       album,
		AlbumArtist: albumArtist,
	})

	if err == nil {
		log.Printf("[%s] NowPlaying: %s by %s\n", api.Name(), title, artist)
	}

	return err
}

package lastfm

import "strconv"

type Args interface {
	Format() map[string]string
}

type ScrobbleArgs struct {
	Artist      string
	Track       string
	Album       string
	AlbumArtist string
	Timestamp   int64
}

func (a ScrobbleArgs) Format() map[string]string {
	return map[string]string{
		"artist":      a.Artist,
		"track":       a.Track,
		"album":       a.Album,
		"albumArtist": a.AlbumArtist,
		"timestamp":   strconv.FormatInt(a.Timestamp, 10),
	}
}

type UpdateNowPlayingArgs struct {
	Artist      string
	Track       string
	Album       string
	AlbumArtist string
}

func (a UpdateNowPlayingArgs) Format() map[string]string {
	return map[string]string{
		"artist":      a.Artist,
		"track":       a.Track,
		"album":       a.Album,
		"albumArtist": a.AlbumArtist,
	}
}

type LoginArgs struct {
	Username string
	Password string
}

func (a LoginArgs) Format() map[string]string {
	return map[string]string{
		"username": a.Username,
		"password": a.Password,
	}
}

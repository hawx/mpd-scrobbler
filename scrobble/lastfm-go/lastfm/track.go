package lastfm

import "strconv"

type trackApi struct {
	uriBase string
	params  *apiParams
}

type ScrobbleArgs struct {
	Artist      string
	Track       string
	Album       string
	AlbumArtist string
	Timestamp   int64
}

func (a ScrobbleArgs) Format() (map[string]string, error) {
	return map[string]string{
		"artist":      a.Artist,
		"track":       a.Track,
		"album":       a.Album,
		"albumArtist": a.AlbumArtist,
		"timestamp":   strconv.FormatInt(a.Timestamp, 10),
	}, nil
}

//track.scrobble
func (api trackApi) Scrobble(args ScrobbleArgs) (result TrackScrobble, err error) {
	err = callPost("track.scrobble", api.uriBase, api.params, args, &result, true)
	return
}

type UpdateNowPlayingArgs struct {
	Artist      string
	Track       string
	Album       string
	AlbumArtist string
}

func (a UpdateNowPlayingArgs) Format() (map[string]string, error) {
	return map[string]string{
		"artist":      a.Artist,
		"track":       a.Track,
		"album":       a.Album,
		"albumArtist": a.AlbumArtist,
	}, nil
}

//track.updateNowPlaying
func (api trackApi) UpdateNowPlaying(args UpdateNowPlayingArgs) (result TrackUpdateNowPlaying, err error) {
	err = callPost("track.updatenowplaying", api.uriBase, api.params, args, &result, true)
	return
}

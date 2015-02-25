package lastfm

type trackApi struct {
	uriBase string
	params  *apiParams
}

//track.scrobble
func (api trackApi) Scrobble(args map[string]interface{}) (result TrackScrobble, err error) {
	defer func() { appendCaller(err, "lastfm.Track.Scrobble") }()
	err = callPost("track.scrobble", api.uriBase, api.params, args, &result, P{
		"indexing": []string{"artist", "track", "timestamp", "album", "context", "streamId", "chosenByUser", "trackNumber", "mbid", "albumArtist", "duration"},
	})
	return
}

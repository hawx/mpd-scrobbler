package lastfm

type trackApi struct {
	uriBase string
	params  *apiParams
}

//track.scrobble
func (api trackApi) Scrobble(args map[string]interface{}) (result TrackScrobble, err error) {
	err = callPost("track.scrobble", api.uriBase, api.params, args, &result, P{
		"indexing": []string{"artist", "track", "timestamp", "album", "context", "streamId", "chosenByUser", "trackNumber", "mbid", "albumArtist", "duration"},
	}, true)
	return
}

//track.updateNowPlaying
func (api trackApi) UpdateNowPlaying(args map[string]interface{}) (result TrackUpdateNowPlaying, err error) {
	err = callPost("track.updatenowplaying", api.uriBase, api.params, args, &result, P{
		"plain": []string{"artist", "track", "album", "trackNumber", "context", "mbid", "duration", "albumArtist"},
	}, true)
	return
}

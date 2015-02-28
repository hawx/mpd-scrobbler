package lastfm

import "encoding/xml"

//track.scrobble
type TrackScrobble struct {
	XMLName   xml.Name `xml:"scrobbles"`
	Accepted  string   `xml:"accepted,attr"`
	Ignored   string   `xml:"ignored,attr"`
	Scrobbles []struct {
		Track struct {
			Corrected string `xml:"corrected,attr"`
			Name      string `xml:",chardata"`
		} `xml:"track"`
		Artist struct {
			Corrected string `xml:"corrected,attr"`
			Name      string `xml:",chardata"`
		} `xml:"artist"`
		Album struct {
			Corrected string `xml:"corrected,attr"`
			Name      string `xml:",chardata"`
		} `xml:"album"`
		AlbumArtist struct {
			Corrected string `xml:"corrected,attr"`
			Name      string `xml:",chardata"`
		} `xml:"albumArtist"`
		TimeStamp      string `xml:"timestamp"`
		IgnoredMessage struct {
			Corrected string `xml:"corrected,attr"`
			Body      string `xml:",chardata"`
		} `xml:"ignoredMessage"`
	} `xml:"scrobble"`
}

//track.updateNowPlaying
type TrackUpdateNowPlaying struct {
	XMLName xml.Name `xml:"nowplaying"`
	Track   struct {
		Corrected string `xml:"corrected,attr"`
		Name      string `xml:",chardata"`
	} `xml:"track"`
	Artist struct {
		Corrected string `xml:"corrected,attr"`
		Name      string `xml:",chardata"`
	} `xml:"artist"`
	Album struct {
		Corrected string `xml:"corrected,attr"`
		Name      string `xml:",chardata"`
	} `xml:"album"`
	AlbumArtist struct {
		Corrected string `xml:"corrected,attr"`
		Name      string `xml:",chardata"`
	} `xml:"albumArtist"`
	IgnoredMessage struct {
		Corrected string `xml:"corrected,attr"`
		Body      string `xml:",chardata"`
	} `xml:"ignoredMessage"`
}

package lastfm

const UriApiSecBase = "https://ws.audioscrobbler.com/2.0/"

type Api struct {
	uriBase string
	params  *apiParams
}

type apiParams struct {
	apikey string
	secret string
	sk     string
}

func New(key, secret, uriBase string) *Api {
	params := apiParams{key, secret, ""}
	if uriBase == "" {
		uriBase = UriApiSecBase
	}

	return &Api{uriBase: uriBase, params: &params}
}

func (api *Api) Scrobble(args ScrobbleArgs) error {
	return api.callPost("track.scrobble", args, nil, true)
}

func (api *Api) UpdateNowPlaying(args UpdateNowPlayingArgs) error {
	return api.callPost("track.updatenowplaying", args, nil, true)
}

func (api *Api) Login(username, password string) error {
	var result AuthGetMobileSession

	if err := api.callPost("auth.getmobilesession", LoginArgs{username, password}, &result, false); err != nil {
		return err
	}
	api.params.sk = result.Key
	return nil
}

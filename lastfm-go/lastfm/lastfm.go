package lastfm

const (
	UriApiSecBase  = "https://ws.audioscrobbler.com/2.0/"
	UriApiBase     = "http://ws.audioscrobbler.com/2.0/"
	UriBrowserBase = "https://www.last.fm/api/auth/"
)

type P map[string]interface{}

type Api struct {
	uriBase string
	params  *apiParams
	Track   *trackApi
}

type apiParams struct {
	apikey    string
	secret    string
	sk        string
	useragent string
}

func New(key, secret, uriBase string) (api *Api) {
	params := apiParams{key, secret, "", ""}
	if uriBase == "" {
		uriBase = UriApiSecBase
	}

	api = &Api{
		uriBase: uriBase,
		params:  &params,
		Track:   &trackApi{uriBase, &params},
	}
	return
}

func (api *Api) SetSession(sessionkey string) {
	api.params.sk = sessionkey
}

func (api Api) GetSessionKey() (sk string) {
	sk = api.params.sk
	return
}

func (api *Api) SetUserAgent(useragent string) {
	api.params.useragent = useragent
}

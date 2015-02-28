package lastfm

//auth.getMobileSession
type AuthGetMobileSession struct {
	Name       string `xml:"name"` //username
	Key        string `xml:"key"`  //session key
	Subscriber bool   `xml:"subscriber"`
}

//Mobile app style
func (api *Api) Login(username, password string) (err error) {
	var result AuthGetMobileSession
	args := P{"username": username, "password": password}
	if err = callPost("auth.getmobilesession", api.uriBase, api.params, args, &result, P{
		"plain": []string{"username", "password"},
	}, false); err != nil {
		return
	}
	api.params.sk = result.Key
	return
}

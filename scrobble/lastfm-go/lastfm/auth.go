package lastfm

//auth.getMobileSession
type AuthGetMobileSession struct {
	Name       string `xml:"name"` //username
	Key        string `xml:"key"`  //session key
	Subscriber bool   `xml:"subscriber"`
}

type LoginArgs struct {
	Username string
	Password string
}

func (a LoginArgs) Format() (map[string]string, error) {
	return map[string]string{
		"username": a.Username,
		"password": a.Password,
	}, nil
}

//Mobile app style
func (api *Api) Login(username, password string) (err error) {
	var result AuthGetMobileSession

	if err = callPost("auth.getmobilesession", api.uriBase, api.params, LoginArgs{username, password}, &result, false); err != nil {
		return
	}
	api.params.sk = result.Key
	return
}

package lastfm

import "encoding/xml"

type Base struct {
	XMLName xml.Name `xml:"lfm"`
	Status  string   `xml:"status,attr"`
	Inner   []byte   `xml:",innerxml"`
}

type ApiError struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:",chardata"`
}

type AuthGetMobileSession struct {
	Name       string `xml:"name"` //username
	Key        string `xml:"key"`  //session key
	Subscriber bool   `xml:"subscriber"`
}

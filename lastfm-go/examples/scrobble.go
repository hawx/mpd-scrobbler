package main

import (
	"../lastfm"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func getTrimmedString(r *bufio.Reader, msg string) (res string) {
	fmt.Print(msg, ":")
	res, _ = r.ReadString('\n')
	res = strings.Trim(res, "\r\n")
	return
}

func main() {
	r := bufio.NewReader(os.Stdin)
	apiKey := getTrimmedString(r, "API KEY")
	apiSecret := getTrimmedString(r, "API SECRET")

	api := lastfm.New(apiKey, apiSecret)

	username := getTrimmedString(r, "Username")
	password := getTrimmedString(r, "Password")

	err := api.Login(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}

	artist := getTrimmedString(r, "Artist")
	track := getTrimmedString(r, "Track")

	p := lastfm.P{"artist": artist, "track": track}
	_, err = api.Track.UpdateNowPlaying(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Updated Now-Playing.")
	start := time.Now().Unix()

	time.Sleep(35 * time.Second)

	p["timestamp"] = start
	_, err = api.Track.Scrobble(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Scrobbled.")
}

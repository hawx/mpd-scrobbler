# mpdrobble

Scrobbler for mpd.

Install with:

``` bash
$ go get github.com/hawx/mpdrobble
```

Then create a config file somewhere, for example
`~/.config/mpdrobble/config.toml`, in the format:

``` toml
[lastfm]
key = "...your lastfm api key..."
secret = "...your lastfm secret..."
username = "...your lastfm username..."
password = "...your lastfm password..."

[trobble]
key = "...your trobble api key..."
secret = "...your trobble secret..."
username = "...your trobble username..."
password = "...your trobble password..."
uri = "...your trobble uri..."
```

As shown multiple sections can be added for other services by also specifying
the uri to the api endpoint. It can then be run with:

``` bash
$ mpdrobble --config ~/.config/mpdrobble/config.toml
...
```

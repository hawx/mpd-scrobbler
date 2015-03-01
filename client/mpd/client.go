// Copyright 2009 The GoMPD Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package mpd provides the client side interface to MPD (Music Player Daemon).
// The protocol reference can be found at http://www.musicpd.org/doc/protocol/index.html
package mpd

import (
	"fmt"
	"net/textproto"
	"strconv"
	"strings"
)

// Client represents a client connection to a MPD server.
type Client struct {
	text *textproto.Conn
}

// Attrs is a set of attributes returned by MPD.
type Attrs map[string]string

// Dial connects to MPD listening on address addr (e.g. "127.0.0.1:6600")
// on network network (e.g. "tcp").
func Dial(network, addr string) (c *Client, err error) {
	text, err := textproto.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	line, err := text.ReadLine()
	if err != nil {
		return nil, err
	}
	if line[0:6] != "OK MPD" {
		return nil, textproto.ProtocolError("no greeting")
	}
	return &Client{text: text}, nil
}

// DialAuthenticated connects to MPD listening on address addr (e.g. "127.0.0.1:6600")
// on network network (e.g. "tcp"). It then authenticates with MPD
// using the plaintext password password if it's not empty.
func DialAuthenticated(network, addr, password string) (c *Client, err error) {
	c, err = Dial(network, addr)
	if err == nil && len(password) > 0 {
		err = c.okCmd("password %s", password)
	}
	return c, err
}

// We are reimplemeting Cmd() and PrintfLine() from textproto here, because
// the original functions append CR-LF to the end of commands. This behavior
// voilates the MPD protocol: Commands must be terminated by '\n'.
func (c *Client) cmd(format string, args ...interface{}) (uint, error) {
	id := c.text.Next()
	c.text.StartRequest(id)
	defer c.text.EndRequest(id)
	if err := c.printfLine(format, args...); err != nil {
		return 0, err
	}
	return id, nil
}

func (c *Client) printfLine(format string, args ...interface{}) error {
	fmt.Fprintf(c.text.W, format, args...)
	c.text.W.WriteByte('\n')
	return c.text.W.Flush()
}

// Close terminates the connection with MPD.
func (c *Client) Close() (err error) {
	if c.text != nil {
		c.printfLine("close")
		err = c.text.Close()
		c.text = nil
	}
	return
}

// Ping sends a no-op message to MPD. It's useful for keeping the connection alive.
func (c *Client) Ping() error {
	return c.okCmd("ping")
}

func (c *Client) readAttrs(terminator string) (attrs Attrs, err error) {
	attrs = make(Attrs)
	for {
		line, err := c.text.ReadLine()
		if err != nil {
			return nil, err
		}
		if line == terminator {
			break
		}
		z := strings.Index(line, ": ")
		if z < 0 {
			return nil, textproto.ProtocolError("can't parse line: " + line)
		}
		key := line[0:z]
		attrs[key] = line[z+2:]
	}
	return
}

type Song struct {
	Title, Artist, Album, AlbumArtist, File string
}

type Pos struct {
	Percent float64
	Seconds int // status time
	Length  int // status time
}

func (c *Client) CurrentSong() (Song, error) {
	s, err := c.cmdReadAttrs("currentsong")
	if err != nil {
		return Song{}, nil
	}

	return Song{s["Title"], s["Artist"], s["Album"], s["AlbumArtist"], s["file"]}, nil
}

func (c *Client) CurrentPos() (pos Pos, playing bool, err error) {
	st, err := c.cmdReadAttrs("status")
	if err != nil {
		return
	}

	playing = true
	if st["volume"] == "-1" {
		playing = false
		return
	}

	parts := strings.Split(st["time"], ":")

	pos.Seconds, err = strconv.Atoi(parts[0])
	if err != nil {
		return
	}
	pos.Length, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}
	pos.Percent = float64(pos.Seconds) / float64(pos.Length) * 100
	return
}

func (c *Client) PlayTime() (int, error) {
	s, err := c.cmdReadAttrs("stats")
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(s["playtime"])
}

func (c *Client) cmdReadAttrs(cmd string) (Attrs, error) {
	id, err := c.cmd(cmd)
	if err != nil {
		return nil, err
	}
	c.text.StartResponse(id)
	defer c.text.EndResponse(id)
	return c.readAttrs("OK")
}

func (c *Client) readOKLine(terminator string) (err error) {
	line, err := c.text.ReadLine()
	if err != nil {
		return
	}
	if line == terminator {
		return nil
	}
	return textproto.ProtocolError("unexpected response: " + line)
}

func (c *Client) okCmd(format string, args ...interface{}) error {
	id, err := c.cmd(format, args...)
	if err != nil {
		return err
	}
	c.text.StartResponse(id)
	defer c.text.EndResponse(id)
	return c.readOKLine("OK")
}

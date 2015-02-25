// Copyright 2013 The GoMPD Authors. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package mpd

import (
	"container/list"
	"errors"
	"strconv"
)

type cmdType uint

const (
	cmdNoReturn cmdType = iota
	cmdAttrReturn
	cmdIDReturn
)

type command struct {
	cmd     string
	promise interface{}
	typeOf  cmdType
}

// CommandList is for batch/mass MPD commands.
// See http://www.musicpd.org/doc/protocol/command_lists.html
// for more details.
type CommandList struct {
	client *Client
	cmdQ   *list.List
}

// PromisedAttrs is a set of promised attributes (to be) returned by MPD.
type PromisedAttrs struct {
	attrs    Attrs
	computed bool
}

func newPromisedAttrs() *PromisedAttrs {
	return &PromisedAttrs{attrs: make(Attrs), computed: false}
}

// PromisedID is a promised identifier (to be) returned by MPD.
type PromisedID int

// Value is a convenience method for ensuring that a promise
// has been computed, returning the Attrs.
func (pa *PromisedAttrs) Value() (Attrs, error) {
	if !pa.computed {
		return nil, errors.New("value has not been computed yet")
	}
	return pa.attrs, nil
}

// Value is a convenience method for ensuring that a promise
// has been computed, returning the ID.
func (pi *PromisedID) Value() (int, error) {
	if *pi == -1 {
		return -1, errors.New("value has not been computed yet")
	}
	return (int)(*pi), nil
}

// BeginCommandList creates a new CommandList structure using
// this connection.
func (c *Client) BeginCommandList() *CommandList {
	return &CommandList{c, list.New()}
}

// Ping sends a no-op message to MPD. It's useful for keeping the connection alive.
func (cl *CommandList) Ping() {
	cl.cmdQ.PushBack(&command{"ping", nil, cmdNoReturn})
}

// CurrentSong returns information about the current song in the playlist.
func (cl *CommandList) CurrentSong() *PromisedAttrs {
	pa := newPromisedAttrs()
	cl.cmdQ.PushBack(&command{"currentsong", pa, cmdAttrReturn})
	return pa
}

// Status returns information about the current status of MPD.
func (cl *CommandList) Status() *PromisedAttrs {
	pa := newPromisedAttrs()
	cl.cmdQ.PushBack(&command{"status", pa, cmdAttrReturn})
	return pa
}

// End executes the command list.
func (cl *CommandList) End() error {

	// Tell MPD to start an OK command list:
	beginID, beginErr := cl.client.cmd("command_list_ok_begin")
	if beginErr != nil {
		return beginErr
	}
	cl.client.text.StartResponse(beginID)
	cl.client.text.EndResponse(beginID)

	// Ensure the queue is cleared regardless.
	defer cl.cmdQ.Init()

	// Issue all of the queued up commands in the list:
	for e := cl.cmdQ.Front(); e != nil; e = e.Next() {
		cmdID, cmdErr := cl.client.cmd(e.Value.(*command).cmd)
		if cmdErr != nil {
			return cmdErr
		}
		cl.client.text.StartResponse(cmdID)
		cl.client.text.EndResponse(cmdID)
	}

	// Tell MPD to end the command list and do the operations.
	endID, endErr := cl.client.cmd("command_list_end")
	if endErr != nil {
		return endErr
	}
	cl.client.text.StartResponse(endID)
	defer cl.client.text.EndResponse(endID)

	// Get the responses back and check for errors:
	for e := cl.cmdQ.Front(); e != nil; e = e.Next() {
		switch e.Value.(*command).typeOf {

		case cmdNoReturn:
			if err := cl.client.readOKLine("list_OK"); err != nil {
				return err
			}

		case cmdAttrReturn:
			a, aErr := cl.client.readAttrs("list_OK")
			if aErr != nil {
				return aErr
			}
			pa := e.Value.(*command).promise.(*PromisedAttrs)
			pa.attrs = a
			pa.computed = true

		case cmdIDReturn:
			a, aErr := cl.client.readAttrs("list_OK")
			if aErr != nil {
				return aErr
			}
			rid, ridErr := strconv.Atoi(a["Id"])
			if ridErr != nil {
				return ridErr
			}
			*(e.Value.(*command).promise.(*PromisedID)) = PromisedID(rid)

		}
	}

	// Finalize the command list with the last OK:
	if cerr := cl.client.readOKLine("OK"); cerr != nil {
		return cerr
	}

	return nil

}

package scrobble

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

var QUEUE_EMPTY = errors.New("Queue empty")

type Track struct {
	Artist      string
	Album       string
	AlbumArtist string
	Title       string
	Timestamp   time.Time
}

type Database interface {
	Queue(name []byte) (Queue, error)
	Close() error
}

type Queue interface {
	Enqueue(Track) error
	Dequeue() (Track, error)
}

type database struct {
	*bolt.DB
}

func Open(path string) (Database, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &database{db}, nil
}

func (this *database) Queue(name []byte) (Queue, error) {
	err := this.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("queue: %s", err)
	}

	return &queue{this.DB, name}, nil
}

type queue struct {
	*bolt.DB
	name []byte
}

func (this *queue) Enqueue(track Track) error {
	return this.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(this.name)
		n, err := b.NextSequence()
		if err != nil {
			return err
		}
		key := strconv.FormatUint(n, 10)
		val, err := json.Marshal(track)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), val)
	})
}

func (this *queue) Dequeue() (track Track, err error) {
	ok := true
	err = this.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(this.name)
		c := b.Cursor()

		key, val := c.First()
		if key == nil {
			ok = false
			return nil
		}

		err := json.Unmarshal(val, &track)
		if err != nil {
			return err
		}
		return b.Delete(key)
	})

	if !ok {
		err = QUEUE_EMPTY
	}
	return
}

package scrobble

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	path := filepath.Join(os.TempDir(), "data_test")
	db, _ := Open(path)
	defer os.Remove(path)

	t1 := Track{"John", "Cool", "John", "Track 01", time.Now().UTC()}
	t2 := Track{"John", "Cool", "John", "Another1", time.Now().UTC()}
	t3 := Track{"Dave", "What", "YesDave", "What2", time.Now().UTC()}

	assert := assert.New(t)

	q, _ := db.Queue([]byte("cool"))
	assert.Nil(q.Enqueue(t1))
	assert.Nil(q.Enqueue(t2))
	assert.Nil(q.Enqueue(t3))

	r1, err1 := q.Dequeue()
	r2, err2 := q.Dequeue()
	r3, err3 := q.Dequeue()
	_, err := q.Dequeue()

	assert.Nil(err1)
	assert.Nil(err2)
	assert.Nil(err3)

	assert.Equal(t1, r1)
	assert.Equal(t2, r2)
	assert.Equal(t3, r3)
	assert.Equal(err, QUEUE_EMPTY)
}

package kv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jialeicui/feedpilot/pkg/store"
)

func TestNewBadger(t *testing.T) {
	var (
		must   = require.New(t)
		dir    = t.TempDir()
		b, err = NewBadger(dir)
	)
	must.NoError(err)
	must.NotNil(b)
}

func TestBadger(t *testing.T) {
	var (
		must   = require.New(t)
		dir    = t.TempDir()
		b, err = NewBadger(dir)
	)
	must.NoError(err)
	must.NotNil(b)

	err = b.Put("key", testStringer("foo"))
	must.NoError(err)

	v, err := b.Get("key")
	must.NoError(err)
	must.Equal("foo", v)

	bWithBucket := b.WithBucket("bucket")
	err = bWithBucket.Put("key", testStringer("bar"))
	must.NoError(err)

	v, err = bWithBucket.Get("key")
	must.NoError(err)
	must.Equal("bar", v)

	// get from the original bucket and it should return the original value
	v, err = b.Get("key")
	must.NoError(err)
	must.Equal("foo", v)

	// test list
	list, err := b.List(0, 0)
	must.NoError(err)
	must.ElementsMatch([]string{"key"}, list)

	// test list with offset and limit
	list, err = b.List(1, 1)
	must.NoError(err)
	must.ElementsMatch([]string{}, list)

	for _, db := range []store.KvStore{b, bWithBucket} {
		// put another key then test list with offset and limit
		err = db.Put("key2", testStringer("bar"))
		must.NoError(err)
		list, err = db.List(1, 1)
		must.NoError(err)
		must.ElementsMatch([]string{"key2"}, list)
		list, err = db.List(0, 0)
		must.NoError(err)
		must.ElementsMatch([]string{"key", "key2"}, list)
		list, err = db.List(0, 1)
		must.NoError(err)
		must.ElementsMatch([]string{"key"}, list)
	}
}

type testStringer string

func (t testStringer) String() string {
	return string(t)
}

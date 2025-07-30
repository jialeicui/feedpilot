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
	// TODO refactor this test
	var (
		must   = require.New(t)
		dir    = t.TempDir()
		b, err = NewBadger(dir)
	)
	must.NoError(err)
	must.NotNil(b)

	err = b.Put("key", newStringer("foo"))
	must.NoError(err)

	v, err := b.Get("key")
	must.NoError(err)
	must.Equal("foo", v)

	bWithBucket := b.WithBucket("bucket")
	err = bWithBucket.Put("key", newStringer("bar"))
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
		err = db.Put("key2", newStringer("bar"))
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

	// delete non-existing key
	err = b.Delete("non-existing-key")
	must.NoError(err)

	// delete existing key
	err = b.Delete("key")
	must.NoError(err)

	_, err = b.Get("key")
	must.Error(err)
}

func TestBadgerSearch(t *testing.T) {
	var (
		must   = require.New(t)
		dir    = t.TempDir()
		b, err = NewBadger(dir)
	)
	must.NoError(err)
	must.NotNil(b)

	// Put some test data
	err = b.Put("user1", newStringer(`{"username":"alice","display_name":"Alice Smith","bio":"Software developer"}`))
	must.NoError(err)
	err = b.Put("user2", newStringer(`{"username":"bob","display_name":"Bob Jones","bio":"Designer and artist"}`))
	must.NoError(err)
	err = b.Put("post1", newStringer(`{"text":"Hello world this is my first post"}`))
	must.NoError(err)
	err = b.Put("post2", newStringer(`{"text":"Another day, another commit"}`))
	must.NoError(err)

	// Test search for "alice" - should find user1
	keys, err := b.Search("alice", 0, 10)
	must.NoError(err)
	must.ElementsMatch([]string{"user1"}, keys)

	// Test search for "developer" - should find user1
	keys, err = b.Search("developer", 0, 10)
	must.NoError(err)
	must.ElementsMatch([]string{"user1"}, keys)

	// Test search for "artist" - should find user2
	keys, err = b.Search("artist", 0, 10)
	must.NoError(err)
	must.ElementsMatch([]string{"user2"}, keys)

	// Test search for "Hello" - should find post1
	keys, err = b.Search("Hello", 0, 10)
	must.NoError(err)
	must.ElementsMatch([]string{"post1"}, keys)

	// Test search for "commit" - should find post2
	keys, err = b.Search("commit", 0, 10)
	must.NoError(err)
	must.ElementsMatch([]string{"post2"}, keys)

	// Test search for "post" - should find post1
	keys, err = b.Search("post", 0, 10)
	must.NoError(err)
	must.ElementsMatch([]string{"post1"}, keys)

	// Test case insensitive search
	keys, err = b.Search("ALICE", 0, 10)
	must.NoError(err)
	must.ElementsMatch([]string{"user1"}, keys)

	// Test search with no results
	keys, err = b.Search("nonexistent", 0, 10)
	must.NoError(err)
	must.Empty(keys)

	// Test search with pagination
	keys, err = b.Search("e", 0, 1) // should find at least one item containing "e"
	must.NoError(err)
	must.Len(keys, 1)
}

type testStringer string

func (t *testStringer) String() (string, error) {
	return string(*t), nil
}

func (t *testStringer) Load(s string) error {
	*t = testStringer(s)
	return nil
}

func newStringer(s string) *testStringer {
	return (*testStringer)(&s)
}

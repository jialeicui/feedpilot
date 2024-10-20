package kv

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"

	"github.com/jialeicui/feedpilot/pkg/store"
)

const (
	defaultBucket = "default"
)

var _ store.KvStore = (*Badger)(nil)

type badgerOptions struct {
	bucket string
}

type BadgerOption func(*badgerOptions)

func WithBucket(bucket string) func(*badgerOptions) {
	return func(opt *badgerOptions) {
		opt.bucket = bucket
	}
}

func applyBadgerOptions(opts []BadgerOption) *badgerOptions {
	opt := &badgerOptions{
		bucket: defaultBucket,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

type Badger struct {
	store  *badger.DB
	bucket string
}

func NewBadger(path string, opt ...BadgerOption) (*Badger, error) {
	badgerOpt := badger.DefaultOptions(path)
	s, err := badger.Open(badgerOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger store: %w", err)
	}
	opts := applyBadgerOptions(opt)
	return &Badger{
		store:  s,
		bucket: opts.bucket,
	}, nil
}

// WithBucket returns a new Badger with the specified bucket.
// it shares the same underlying store.
func (b *Badger) WithBucket(bucket string) *Badger {
	ret := *b
	ret.bucket = bucket
	return &ret
}

func (b *Badger) Put(key string, value store.Stringer) error {
	key = fmt.Sprintf("%s/%s", b.bucket, key)
	return b.store.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value.String()))
	})
}

func (b *Badger) Get(key string) (string, error) {
	key = fmt.Sprintf("%s/%s", b.bucket, key)
	var value string
	err := b.store.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
	})
	if err != nil {
		return "", err
	}
	return value, nil
}

// List returns a list of keys with offset and limit
func (b *Badger) List(offset, limit int) ([]string, error) {
	var (
		keys   []string
		cursor int
	)
	err := b.store.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		firstKey := []byte(fmt.Sprintf("%s/", b.bucket))
		for it.Seek(firstKey); it.ValidForPrefix(firstKey); it.Next() {
			cursor++
			if offset > 0 && cursor <= offset {
				continue
			}

			item := it.Item()
			k := item.Key()
			// remove the bucket prefix
			k = k[len(b.bucket)+1:]
			keys = append(keys, string(k))
			if limit > 0 && len(keys) >= limit {
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return keys, nil
}

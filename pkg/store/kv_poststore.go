package store

import (
	"github.com/jialeicui/feedpilot/pkg/meta"
)

var _ PostStore = (*postStore)(nil)

type postStore struct {
	store KvStore
}

func (p *postStore) Put(id meta.PostID, post *meta.Post) error {
	return p.store.Put(id.String(), post)
}

func (p *postStore) Get(id meta.PostID) (*meta.Post, error) {
	val, err := p.store.Get(id.String())
	if err != nil {
		return nil, err
	}
	var post = new(meta.Post)
	if err := post.Load(val); err != nil {
		return nil, err
	}
	return post, nil
}

func (p *postStore) Delete(id meta.PostID) error {
	return p.store.Delete(id.String())
}

func (p *postStore) List(offset, limit int) ([]*meta.Post, error) {
	keys, err := p.store.List(offset, limit)
	if err != nil {
		return nil, err
	}
	posts := make([]*meta.Post, 0, len(keys))
	for _, key := range keys {
		val, err := p.store.Get(key)
		if err != nil {
			return nil, err
		}
		var post = new(meta.Post)
		if err := post.Load(val); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *postStore) Search(query string, offset, limit int) ([]*meta.Post, error) {
	keys, err := p.store.Search(query, offset, limit)
	if err != nil {
		return nil, err
	}
	posts := make([]*meta.Post, 0, len(keys))
	for _, key := range keys {
		val, err := p.store.Get(key)
		if err != nil {
			return nil, err
		}
		var post = new(meta.Post)
		if err := post.Load(val); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func NewPostStore(kv KvStore) PostStore {
	return &postStore{
		store: kv,
	}
}
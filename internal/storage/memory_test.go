package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// verify that the roundrip create -> get -> delete -> get behaves as expected
func TestMemoryRoundtrip(t *testing.T) {
	store, err := NewMemoryStorage()
	require.NoError(t, err, "shouldn't fail initing the store")

	const target = "https://google.com"

	slug, err := store.CreateLink(target, nil)
	require.NoError(t, err, "creating a new link shouldn't error on this test")
	assert.NotEqual(t, "", slug, "the slug should never be empty")

	links := store.AllLinks()
	assert.Equal(t, 1, len(links), "there should be only one link on the database")

	link, err := store.GetLink(slug)
	require.NoError(t, err, "this slug was returned by the previous create call; it should be on db")

	assert.Equal(t, target, link.Target, "the slug had a link on the db, but not for the correct url?")

	err = store.DeleteLink(slug)
	require.NoError(t, err, "the in-memory db doesn't error on delete; and the slug should exist")

	links = store.AllLinks()
	assert.Equal(t, 0, len(links), "db should be empty now")

	_, err = store.GetTarget(slug)
	require.Errorf(t, err, "we deleted this slug, it should fail when trying to fetch it")
}

func TestMemoryHitRegistering(t *testing.T) {
	store, err := NewMemoryStorage()
	require.NoError(t, err, "shouldn't fail initing the store")

	const target = "https://google.com"

	slug, err := store.CreateLink(target, nil)
	require.NoError(t, err, "creating a new link shouldn't error on this test")
	assert.NotEqual(t, "", slug, "the slug should never be empty")

	store.RegisterHit(slug)
	store.RegisterHit(slug)

	link, err := store.GetLink(slug)
	require.NoError(t, err, "this slug was returned by the previous create call; it should be on db")

	assert.Equal(t, target, link.Target, "the slug had a link on the db, but not for the correct url?")
	assert.EqualValues(t, 2, link.Hits, "this link should have been visited twice")
}

type staticslugger struct{}

func (s *staticslugger) Random() (string, error) {
	return "test", nil
}

func TestMemorySlugGenRecursion(t *testing.T) {
	store, err := NewMemoryStorage(WithSlugGenerator(&staticslugger{}))
	require.NoError(t, err, "shouldn't fail initing the store")

	const target = "https://google.com"

	slug, err := store.CreateLink(target, nil)
	require.NoError(t, err, "creating a new link the first time shouldn't error")

	assert.Equal(t, slug, "test", "the slug generated is not matching the static slug used on this test")

	_, err = store.CreateLink(target, nil)
	require.Error(t, err, "this time it should error, as the slug generator always returned the same value")
	assert.Contains(t, err.Error(), "generate a unique slug")
}

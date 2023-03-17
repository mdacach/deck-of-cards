package deck

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoreAddDeck(t *testing.T) {
	store := NewStore()

	deck := NewStandardDeck()
	err := store.Add(&deck)
	require.NoError(t, err)

	retrievedDeck, err := store.Get(deck.ID)
	require.NoError(t, err)
	assert.Equal(t, deck, *retrievedDeck)
}

func TestStoreAddMultipleDecks(t *testing.T) {
	store := NewStore()

	decks := []Deck{NewStandardDeck(), NewStandardDeck(), NewStandardDeck()}
	for _, deck := range decks {
		err := store.Add(&deck)
		require.NoError(t, err)
	}

	// All the decks have been correctly inserted.
	for _, deck := range decks {
		_, err := store.Get(deck.ID)
		require.NoError(t, err)
	}
}

func TestStoreGetNonExistentDeck(t *testing.T) {
	store := NewStore()

	nonExistentID := uuid.New()
	_, err := store.Get(nonExistentID)
	assert.Error(t, err)
}

func TestStoreRemoveDeck(t *testing.T) {
	store := NewStore()

	deck := NewStandardDeck()
	err := store.Add(&deck)
	require.NoError(t, err)

	err = store.Remove(deck.ID)
	require.NoError(t, err)

	_, err = store.Get(deck.ID)
	assert.Error(t, err)
}

func TestStoreRemoveNonExistentDeck(t *testing.T) {
	store := NewStore()

	nonExistentID := uuid.New()
	err := store.Remove(nonExistentID)
	assert.Error(t, err)
}

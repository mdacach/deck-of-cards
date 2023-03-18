package deck

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

// Store manages a collection of decks in a thread-safe manner. Decks are accessed by their ID (uuid).
type Store struct {
	decks map[uuid.UUID]*Deck
	// Maps are not safe for concurrent access, so we will synchronize with a Mutex.
	// It is OK to have concurrent reads though, se we use RWMutex.
	mu sync.RWMutex
}

// NewStore creates and returns a new instance of the Store type.
func NewStore() *Store {
	return &Store{
		decks: make(map[uuid.UUID]*Deck),
	}
}

// Add adds a new deck to the store. It returns an error if a deck with the same ID already exists in the store.
func (s *Store) Add(deck *Deck) error {
	s.mu.Lock()
	// From Go 1.14 we don't need to worry about performance of defer here.
	// Source: https: //www.reddit.com/r/golang/comments/fdy6sb/comment/fjl6d37/?utm_source=share&utm_medium=web2x&context=3
	defer s.mu.Unlock()

	// It is better to return an error here then to overwrite a deck. The overwritten deck may have been used.
	// Or overwriting decks could be a potential attack.
	if _, exists := s.decks[deck.ID]; exists {
		return errors.New("deck ID already exists in the store")
	}

	s.decks[deck.ID] = deck
	return nil
}

// Get retrieves a deck from the store by its ID. It returns the deck and nil if the deck is found, or nil and an error if the deck is not found.
func (s *Store) Get(id uuid.UUID) (*Deck, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if deck, ok := s.decks[id]; ok {
		return deck, nil
	}
	return nil, errors.New("deck not found")
}

// Remove removes a deck from the store by its ID. It returns an error if the deck is not found.
func (s *Store) Remove(deckID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.decks[deckID]
	if !exists {
		return errors.New("deck not found")
	}

	delete(s.decks, deckID)
	return nil
}

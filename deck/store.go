package deck

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type Store struct {
	decks map[uuid.UUID]*Deck
	// Maps are not safe for concurrent access, so we will synchronize with a Mutex.
	// It is OK to have concurrent reads though, se we will RWMutex.
	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		decks: make(map[uuid.UUID]*Deck),
	}
}

func (s *Store) Add(deck *Deck) error {
	s.mu.Lock()
	// From Go 1.14 we don't need to worry about performance of defer here.
	// Source: https: //www.reddit.com/r/golang/comments/fdy6sb/comment/fjl6d37/?utm_source=share&utm_medium=web2x&context=3
	defer s.mu.Unlock()

	if _, exists := s.decks[deck.ID]; exists {
		return errors.New("deck ID already exists in the store")
	}

	s.decks[deck.ID] = deck
	return nil
}

func (s *Store) Get(id uuid.UUID) (*Deck, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if deck, ok := s.decks[id]; ok {
		return deck, nil
	}
	return nil, errors.New("deck not found")
}

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

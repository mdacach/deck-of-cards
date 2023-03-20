# Design Document: Deck of Cards API

## Overview

The Deck of Cards API is a web service that provides an interface to interact with a deck of playing cards. The API
allows users to create standard or custom decks, shuffle the decks, draw cards, and manage multiple decks in a store.
The main purpose of this project is to offer a flexible, easy-to-use RESTful API for card game applications.

## Requirements

### Functional Requirements

1. Create a **full** (standard) deck of 52 playing cards, with an option to shuffle the deck.
2. Create a **partial** deck with a specific set of cards, with an option to shuffle the deck.
3. Retrieve information about (**open**) a deck, including the deck ID, remaining cards, and whether the deck is
   shuffled.
4. **Draw** a specified number of cards from a deck.
5. Store and manage multiple decks in a deck store.

### Non-Functional Requirements

- The code should be well tested.
- The code should be easy to modify, maintain and extend.
- The code should be easy to understand.
- The code should comply with best Go practices.

## Architecture

The Deck of Cards API is a web service built using the Go programming language and the Gin web framework. The API
consists of three main packages: `card`, `deck`, and `api`. The `card` package defines the `Card`, `Rank`, and `Suit`
types, while the
`deck` package provides the `Deck` type and deck-related operations. The `api` package handles the RESTful endpoints and
request/response handling.

### Package: card

The `card` package defines the `Card`, `Rank`, and `Suit` types. These types represent the basic elements of a playing
card and
provide methods to parse and validate cards.

### Package: deck

The `deck` package defines the `Deck` type and associated operations such as creating standard and partial decks,
shuffling,
drawing cards, and managing remaining cards in the deck. The package also includes the `Store` type, which allows for
the
in-memory management of multiple decks using a map and mutex for concurrent access control.

### Package: api

The `api` package handles the RESTful endpoints and request/response handling using the Gin web framework. It provides
the
following endpoints:

1. `POST /deck/new`: Create a new deck (full or partial) with optional shuffling.
2. `GET /deck/:deck_id`: Retrieve information about (open) a deck.
3. `GET /deck/:deck_id/draw`: Draw a specified number of cards from a deck.

The package also defines the required request and response structures for each endpoint.

## Use Cases

1. A user creates a standard deck of cards and shuffles it:

   ```console
   POST /deck/new?shuffled=true
   ```

   By default the deck is un-shuffled.

2. A user creates a custom deck with specific cards and shuffles it:

   ```console
   POST /deck/new?cards=AS,KD,2H,5C&shuffled=true
   ```

3. A user opens a deck:

   ```console
   GET /deck/123e4567-e89b-12d3-a456-426655440000
   ```

4. A user draws three cards from a deck:

   ```console
   GET /deck/123e4567-e89b-12d3-a456-426655440000/draw?count=3
   ```
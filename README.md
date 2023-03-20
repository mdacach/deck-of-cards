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
3. `POST /deck/:deck_id/draw`: Draw a specified number of cards from a deck.

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
   POST /deck/123e4567-e89b-12d3-a456-426655440000/draw?count=3
   ```

## Example Usage

Note that the code here will not work on your machine because the uuid of your generated
deck will be different. In order to follow along, you will need to update the deck_id between requests.

### Start the Server

The server runs on port :8080.

```shell
make server
```

or

```shell
go run .
```

if no `make` available.

1. ### Partial Shuffled Deck
   #### Create a partial shuffled Deck
   ```shell
   curl --request POST \
     --url 'http://localhost:8080/deck/new?cards=AS,KD,AC,2C,KH&shuffled=true'
   ```

   <details>
   <summary>Example Output</summary>

   ```json
   {
      "deck_id":"31ef40c2-5825-491c-b5c6-68e385717427",
      "shuffled":true,
      "remaining":5
   }
   ```

   </details>

   #### Open the deck
   ```shell
   curl --request GET \
     --url 'http://localhost:8080/deck/31ef40c2-5825-491c-b5c6-68e385717427'
   ```

   <details>
   <summary>Example Output</summary>

   ```json
   {
      "deck_id":"31ef40c2-5825-491c-b5c6-68e385717427",
      "shuffled":true,
      "remaining":5,
      "cards": [
         {
            "value":"KING",
            "suit":"DIAMONDS",
            "code":"KD"},
         {
            "value":"KING",
            "suit":"HEARTS",
            "code":"KH"},
         {
            "value":"ACE",
            "suit":"SPADES",
            "code":"AS"},
         {
            "value":"TWO",
            "suit":"CLUBS",
            "code":"2C"},
         {
            "value":"ACE",
            "suit":"CLUBS",
            "code":"AC"
         }
      ]
   }
   ```

   </details>

   #### Draw three cards
   ```shell
   curl --request POST \
     --url 'http://localhost:8080/deck/31ef40c2-5825-491c-b5c6-68e385717427/draw?count=3'
   ```

   <details>
   <summary>Example Output</summary>

   ```json
   {
      "cards": [
         {
            "value":"KING",
            "suit":"DIAMONDS",
            "code":"KD"
         },
         {
            "value":"KING",
            "suit":"HEARTS",
            "code":"KH"
         },
         {
            "value":"ACE",
            "suit":"SPADES",
            "code":"AS"
         }
      ]
   }
   ```

   </details>
2. ### Poker Hand
   #### Create standard deck
   ```shell
   curl --request POST \
     --url 'http://localhost:8080/deck/new?shuffled=true'
   ```

   <details>
   <summary>Example Output</summary>

   ```json
   {
      "deck_id":"9c604247-f9f2-4bd3-8354-00c30463ec37",
      "shuffled":true,
      "remaining":52
   }
   ```

   </details>

   #### Draw hand
   ```shell
   curl --request POST \
     --url 'http://localhost:8080/deck/9c604247-f9f2-4bd3-8354-00c30463ec37/draw?count=2'
   ```

   <details>
   <summary>Example Output</summary>
   Ts6c is not a good hand... There are some chances at a straight though, let's see the flop.

   ```json
   {
      "cards": [
         {
            "value":"TEN",
            "suit":"SPADES",
            "code":"TS"
         },
         {
            "value":"SIX",
            "suit":"CLUBS",
            "code":"6C"
         }
      ]
   }
   ```

   </details>

   #### Draw the flop
   ```shell
   curl --request POST \
     --url 'http://localhost:8080/deck/9c604247-f9f2-4bd3-8354-00c30463ec37/draw?count=3'
   ```

   <details>
   <summary>Example Output</summary>
   We do have a shot at a straight (if any 7 comes)...  
   But it's not really enough to continue.  
   Let's fold here.

   ```json
   {
      "cards": [
         {
            "value":"EIGHT",
            "suit":"CLUBS",
            "code":"8C"
         },
         {
            "value":"NINE",
            "suit":"CLUBS",
            "code":"9C"
         },
         {
            "value":"KING",
            "suit":"DIAMONDS",
            "code":"KD"
         }
      ]
   }
   ```

</details>

## Possible Improvements

### Choose port to run

Currently, the server always run on port :8080. We could add a flag or an environment variable to set
which port to run.

### Data Persistence

We store the decks in-memory, through a Map + Mutex.
This means that if the server shuts down, every stored deck is lost.
We could implement data persistence by using a database.
A NoSQL database could be a good fit, as the data is not tabular and we won't need to perform
complex queries.

### Require authorization for opening decks

Opening the deck exposes all the cards, in-order. If a user was able to inspect the deck,
they would know which cards to expect when drawing.

### Better logging

Gin already provides some logging, but we are not logging anything else.

### Continuous Integration

For a bigger team, it would be nice to automatically reject a code change if tests do not pass.

### Containerize

We could make the project easier to run and develop by using container solutions such as Docker.

### Property-based testing

Some tests (such as CardFromString) could be improved by
using [Property-based testing](https://earthly.dev/blog/property-based-testing/).

### Idempotency of create

User may try to create the same deck multiple times (maybe they click at the create button too many times).
We would treat each request as different, and create multiple decks.

### Better error handling

Error handling could be improved by providing custom error types and better error messages.

### Less memory usage

In order to better reflect the JSON responses, we store all the cards individually inside a deck.
The `deck` type contain an array of cards (with rank, suit and code). This is not necessary, as we only
have a fixed amount (52) different cards, and they do not change internally. We could improve the memory usage by
referring
to the cards, instead of storing the cards. (e.g. store a []int inside Deck, where each integer corresponds to a card by
index
(from 1 to 52)).

### Type-Driven Development

Let's consider `Rank` for example. Rank is a length-one string ("A" for Ace, "K" for King).
But there's no validation for creating a Rank. It is allowed to create a Rank("Some String"), and the compiler will not
complain.
[Newtype](https://www.reddit.com/r/golang/comments/kmj640/newtypes_constructing_and_validation/?utm_source=share&utm_medium=web2x&context=3)
pattern could be useful for handling this.
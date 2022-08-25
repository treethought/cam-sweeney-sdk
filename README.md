# The One API SDK

Golang SDK for [The One API](https://the-one-api.dev/documentation#3)

## Getting started

Add the SDK as dependency

```
go install github.com/treethought/cam-sweeney-sdk@latest
```

Obtain an APIKey from https://the-one-api.dev/sign-up and set the ONE_API_KEY environment variable

```
export ONE_API_KEY=<MY_API_KEY>
```

## Usage

Create a client

```go
import (
    "fmt"
    "os"
    "github.com/treethought/cam-sweeney-sdk"
    )

func main() {
    apiKey := os.GetEnv("ONE_API_KEY")
    client := NewOneAPIClient(apiKey)
}

```

The client struct provides methods to interface with the Books, Movies, Characters, Quotes, and Chapters resources.

### Books

The `Books()` method provides an interface to list and get books.

```go
import (
    "fmt"
    "os"
    "github.com/treethought/cam-sweeney-sdk"
    )

func main() {
    apiKey := os.GetEnv("ONE_API_KEY")
    client := NewOneAPIClient(apiKey)

    // List all available books
    books, err := client.Books().List()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(books)

    // get a specific book
    book, err := client.Books().Get("5cf5805fb53e011a64671582")
    if err != nil {
        log.Fatal(err)
    }
    // get all chapters of a book
    book, err := client.Books().GetChapters("5cf5805fb53e011a64671582")
    if err != nil {
        log.Fatal(err)
    }

}

```

### Movies

The `Movies()` method provides an interface to list and get movies.

```go
// List all available movies
movies, err := client.Movies().List()
if err != nil {
    log.Fatal(err)
}
fmt.Println(movies)

// get a specific movie
book, err := client.Movies().Get("5cd95395de30eff6ebccde5c")
if err != nil {
    log.Fatal(err)
}
// get all quotes of a movie
book, err := client.Movies().GetQuotes("5cd95395de30eff6ebccde5c")
if err != nil {
    log.Fatal(err)
}

```

### Characters

The `Characters()` method provides an interface to list and get characters.

```go
// List all available characters
chapters, err := client.Characters().List()
if err != nil {
    log.Fatal(err)
}
fmt.Println(chapters)

// get a specific chapter
chapter, err := client.Characters().Get("5cd99d4bde30eff6ebccfd0d")
if err != nil {
    log.Fatal(err)
}

// get all quotes of a given character
chapter, err := client.Characters().Get("5cd99d4bde30eff6ebccfd0d")
if err != nil {
    log.Fatal(err)
}

// get all quotes spoken by a character
chapter, err := client.Characters().GetQuotes("5cd99d4bde30eff6ebccfd0d")
if err != nil {
    log.Fatal(err)
}


```

### Chapters

The `Chapters()` method provides an interface to list and get chapters.

```go
// List all available chapters
chapters, err := client.Chapters().List()
if err != nil {
    log.Fatal(err)
}
fmt.Println(chapters)

// get a specific chapter
chapter, err := client.Chapters().Get("6091b6d6d58360f988133b8b")
if err != nil {
    log.Fatal(err)
}

```

### Quotes

The `Quotes)` method provides an interface to list and get quotes.

```go
// List all available quotes
quotes, err := client.Quotes().List()
if err != nil {
    log.Fatal(err)
}
fmt.Println(quotes)

// get a specific quote
quote, err := client.Quotes().Get("5cd96e05de30eff6ebcce7e9")
if err != nil {
    log.Fatal(err)
}

```

## Testing

To test the SDK:

First clone the repo

```

git clone https://github.com/treethought/cam-sweeney-sdk.git
cd cam-sweeney-sdk

```

Run the tests

```

go test -v ./sdk/...
```

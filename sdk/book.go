package sdk

import "fmt"

type booksResponse struct {
	Docs []Book
}

type Book struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name,omitempty"`
}

type BooksClient struct {
	c OneAPIClient
}

type BookClient struct {
	c OneAPIClient
}

// ListBooks returns a list of all "Lord of the Rings" books
func (b BooksClient) List() ([]Book, error) {
	resp := booksResponse{}
	err := b.c.doRequestInto("/book", &resp)
	if err != nil {
		return nil, err
	}
	return resp.Docs, err
}

// Get a book by it's ID
func (b BooksClient) Get(id string) (Book, error) {
	path := fmt.Sprintf("/book/%s", id)
	resp := booksResponse{}
	err := b.c.doRequestInto(path, &resp)
	if err != nil {
		return Book{}, err
	}
	return resp.Docs[0], err
}

func (b BooksClient) GetChapters(bookId string) ([]Chapter, error) {
	path := fmt.Sprintf("/book/%s/chapter", bookId)
	resp := chapterResponse{}
	err := b.c.doRequestInto(path, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Docs, err
}

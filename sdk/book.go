package sdk

type Book struct {
	ID   string
	name string
}

type ListBooksResponse struct {
	docs   []Book
	total  int
	limit  int
	offset int
	page   int
	pages  int
}

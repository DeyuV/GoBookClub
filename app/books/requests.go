package books

type SuggestRequest struct {
	Term string
}

type BookByAuthorOrTitleRequest struct {
	Title  string
	Author string
}

type BookRequest struct {
	Title  string
	Author string
	Year   uint
}

type BookListRequest struct {
	UserID uint
	BookID uint
	Title  string
	Author string
	Year   uint
}

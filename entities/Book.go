package entities

type Book struct {
	BookID        int     `json:"bookID"`
	AuthorID      int     `json:"authorID"`
	Title         string  `json:"title"`
	Publication   string  `json:"publication"`
	PublishedDate string  `json:"publishedDate"`
	Author        *Author `json:",omitempty"`
}

package entities

type Author struct {
	AuthorID  int    `json:"authorID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	DOB       string `json:"DOB"`
	PenName   string `json:"penName"`
}

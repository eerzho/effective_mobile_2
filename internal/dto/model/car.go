package model

type Car struct {
	ID      uint    `json:"id"`
	RegNum  string  `json:"regNum"`
	Mark    string  `json:"mark"`
	Model   string  `json:"model"`
	Year    int     `json:"year"`
	OwnerID uint    `json:"ownerID"`
	Owner   *People `json:"owner,omitempty"`
}

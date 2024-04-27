package model

type Car struct {
	ID      uint    `json:"id"`
	RegNum  string  `json:"reg_num"`
	Mark    string  `json:"mark"`
	Model   string  `json:"model"`
	Year    int     `json:"year"`
	OwnerID uint    `json:"owner_id"`
	Owner   *People `json:"owner,omitempty"`
}

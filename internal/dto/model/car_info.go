package model

type CarInfo struct {
	Mark  string  `json:"mark"`
	Model string  `json:"model"`
	Year  *int    `json:"year"`
	Owner *People `json:"owner,omitempty"`
}

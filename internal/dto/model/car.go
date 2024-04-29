package model

type Car struct {
	ID      uint   `json:"id"`
	RegNum  string `json:"regNum"`
	OwnerID uint   `json:"ownerID"`
	CarInfo
}

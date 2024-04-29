package model

type People struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic"`
}

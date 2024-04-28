package command

type CarIndex struct {
	RegNum *string `schema:"regNum"`
	Mark   *string `schema:"mark"`
	Model  *string `schema:"model"`
	Year   *int    `schema:"year"`
	Page   *int    `schema:"page"`
	Count  *int    `schema:"count"`
}

type CarStore struct {
	RegNums []string `schema:"regNums" validate:"required"`
}

type CarUpdate struct {
	ID     int
	RegNum *string `json:"regNum" validate:"omitempty,ne="`
	Mark   *string `json:"mark" validate:"omitempty,ne="`
	Model  *string `json:"model" validate:"omitempty,ne="`
	Year   *int    `json:"year" validate:"omitempty,gte=1886,lte=2023"`
}

type CarDelete struct {
	ID int
}

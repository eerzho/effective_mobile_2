package request

type CarIndex struct {
	RegNum       *string `schema:"regNum"`
	Mark         *string `schema:"mark"`
	Model        *string `schema:"model"`
	Year         *int    `schema:"year"`
	OwnerName    *string `schema:"ownerName"`
	OwnerSurname *string `schema:"ownerSurname"`
	Order        *string `schema:"order"`
	Page         *int    `schema:"page"`
	Count        *int    `schema:"count"`
}

type CarStore struct {
	RegNums []string `json:"regNums" validate:"required,min=1,dive,required"`
}

type CarUpdate struct {
	RegNum *string `json:"regNum" validate:"omitempty,ne="`
	Mark   *string `json:"mark" validate:"omitempty,ne="`
	Model  *string `json:"model" validate:"omitempty,ne="`
	Year   *int    `json:"year" validate:"omitempty,gte=1886,lte=2023"`
}

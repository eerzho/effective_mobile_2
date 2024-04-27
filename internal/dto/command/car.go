package command

type CarIndex struct {
	Mark         *string `schema:"mark"`
	Model        *string `schema:"model"`
	Year         *int    `schema:"year"`
	OwnerSurname *string `schema:"ownerSurname"`
	Page         *int    `schema:"page"`
	Count        *int    `schema:"count"`
}

package command

type CarIndex struct {
	RegNum       *string
	Mark         *string
	Model        *string
	Year         *int
	OwnerName    *string
	OwnerSurname *string
	Order        *string
	Page         *int
	Count        *int
}

type CarStore struct {
	RegNums []string
}

type CarUpdate struct {
	ID     int
	RegNum *string
	Mark   *string
	Model  *string
	Year   *int
}

type CarDelete struct {
	ID int
}

package query

type CarList struct {
	RegNum *string
	Mark   *string
	Model  *string
	Year   *int
	Page   int
	Count  int
}

type CarCreate struct {
	RegNum  string
	Mark    string
	Model   string
	Year    *int
	OwnerID uint
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

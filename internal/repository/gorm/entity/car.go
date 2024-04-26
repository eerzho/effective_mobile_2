package entity

type Car struct {
	ID      uint   `gorm:"primary_key"`
	RegNum  string `gorm:"unique;not null"`
	Mark    string `gorm:"type:varchar(100)"`
	Model   string `gorm:"type:varchar(100)"`
	Year    int
	Owner   People `gorm:"foreignKey:OwnerID"`
	OwnerID uint
}

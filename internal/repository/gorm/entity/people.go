package entity

type People struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `gorm:"type:varchar(100)"`
	Surname    string `gorm:"type:varchar(100)"`
	Patronymic string `gorm:"type:varchar(100);default:null"`
}

package car

import (
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/repository/gorm/people"
)

type Car struct {
	ID      uint   `gorm:"primary_key"`
	RegNum  string `gorm:"unique;not null"`
	Mark    string `gorm:"type:varchar(100)"`
	Model   string `gorm:"type:varchar(100)"`
	Year    int
	Owner   people.People `gorm:"foreignKey:OwnerID"`
	OwnerID uint
}

func ToModel(entity Car) model.Car {
	carInfo := model.CarInfo{Mark: entity.Mark, Model: entity.Model}
	if entity.Year != 0 {
		carInfo.Year = &entity.Year
	}

	car := model.Car{
		ID:      entity.ID,
		RegNum:  entity.RegNum,
		OwnerID: entity.OwnerID,
		CarInfo: carInfo,
	}

	if entity.Owner.ID != 0 {
		owner := people.ToModel(entity.Owner)
		car.Owner = &owner
	}

	return car
}

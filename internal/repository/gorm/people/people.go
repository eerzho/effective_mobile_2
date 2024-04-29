package people

import "effective_mobile_2/internal/dto/model"

type People struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `gorm:"type:varchar(100)"`
	Surname    string `gorm:"type:varchar(100)"`
	Patronymic string `gorm:"type:varchar(100);default:null"`
}

func ToModel(entity People) model.People {
	people := model.People{
		ID:      entity.ID,
		Name:    entity.Name,
		Surname: entity.Surname,
	}

	if entity.Patronymic != "" {
		people.Patronymic = &entity.Patronymic
	}

	return people
}

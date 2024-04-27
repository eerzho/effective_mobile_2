package poople

import "effective_mobile_2/internal/dto/model"

type People struct {
	ID         uint   `gorm:"primary_key"`
	Name       string `gorm:"type:varchar(100)"`
	Surname    string `gorm:"type:varchar(100)"`
	Patronymic string `gorm:"type:varchar(100);default:null"`
}

func ToModel(entity People) model.People {
	return model.People{
		ID:         entity.ID,
		Name:       entity.Name,
		Surname:    entity.Surname,
		Patronymic: entity.Patronymic,
	}
}

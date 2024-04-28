package car

import (
	"log/slog"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/query"
	"effective_mobile_2/internal/repository/gorm/poople"
	"gorm.io/gorm"
)

type Car struct {
	ID      uint   `gorm:"primary_key"`
	RegNum  string `gorm:"unique;not null"`
	Mark    string `gorm:"type:varchar(100)"`
	Model   string `gorm:"type:varchar(100)"`
	Year    int
	Owner   poople.People `gorm:"foreignKey:OwnerID"`
	OwnerID uint
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(qry *query.CarList) (*[]model.Car, error) {
	const op = "repository.gorm.car.List"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("repository starting", slog.Any("qry", qry))

	log.Info("searching cars")
	builder := r.db.Model(&Car{})

	if qry.RegNum != nil {
		builder.Where("reg_num = ?", *qry.RegNum)
	}
	if qry.Mark != nil {
		builder = builder.Where("mark LIKE ?", "%"+*qry.Mark+"%")
	}
	if qry.Model != nil {
		builder = builder.Where("model LIKE ?", "%"+*qry.Model+"%")
	}
	if qry.Year != nil {
		builder = builder.Where("year = ?", qry.Year)
	}

	builder.Preload("Owner")
	builder = builder.Limit(qry.Count).Offset((qry.Page - 1) * qry.Count)

	var entities []Car
	result := builder.Find(&entities)
	if result.Error != nil {
		log.Error("failed to search", slog.String("error", result.Error.Error()))
		return nil, result.Error
	}
	log.Debug("searched cars", slog.Any("entities", entities))

	log.Info("adapting entity to model")
	cars := make([]model.Car, len(entities))
	for i, entity := range entities {
		cars[i] = ToModel(entity)
	}
	log.Debug("adapted entity to model", slog.Any("cars", cars))

	return &cars, nil
}

func (r *Repository) Update(qry *query.CarUpdate) (*model.Car, error) {
	const op = "repository.car.Update"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("repository starting", slog.Any("qry", qry))

	log.Info("finding car")
	var car Car
	result := r.db.First(&car, qry.ID)
	if result.Error != nil {
		log.Error("failed to find car", slog.String("error", result.Error.Error()))
		return nil, result.Error
	}
	log.Debug("found car", slog.Any("car", car))

	if qry.RegNum != nil {
		car.RegNum = *qry.RegNum
	}
	if qry.Mark != nil {
		car.Mark = *qry.Mark
	}
	if qry.Model != nil {
		car.Model = *qry.Model
	}
	if qry.Year != nil {
		car.Year = *qry.Year
	}

	log.Info("updating", slog.Any("car", car))
	result = r.db.Save(&car)
	if result.Error != nil {
		log.Error("failed to update car", slog.String("error", result.Error.Error()))
		return nil, result.Error
	}
	log.Debug("updated car", slog.Any("car", car))

	log.Info("adapting entity to model")
	m := ToModel(car)
	log.Debug("adapted entity to model", slog.Any("m", m))

	return &m, nil
}

func (r *Repository) Delete(qry *query.CarDelete) error {
	const op = "repository.car.Delete"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("repository starting", slog.Any("qry", qry))

	log.Info("deleting car")
	result := r.db.Delete(&Car{}, qry.ID)
	if result.Error != nil {
		log.Error("failed to delete car", slog.String("error", result.Error.Error()))
		return result.Error
	}

	return nil
}

func ToModel(entity Car) model.Car {
	car := model.Car{
		ID:      entity.ID,
		RegNum:  entity.RegNum,
		Model:   entity.Model,
		Mark:    entity.Mark,
		Year:    entity.Year,
		OwnerID: entity.OwnerID,
	}
	if entity.Owner.ID != 0 {
		owner := poople.ToModel(entity.Owner)
		car.Owner = &owner
	}

	return car
}

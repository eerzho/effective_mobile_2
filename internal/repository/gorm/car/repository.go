package car

import (
	"errors"
	"fmt"
	"log/slog"

	"effective_mobile_2/internal/app_error"
	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/query"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(qry *query.CarList) (*[]model.Car, error) {
	const op = "repository.gorm.car.List"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("qry", qry),
	)

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
	if qry.OwnerName != nil {
		builder.Joins("JOIN peoples ON peoples.id = cars.owner_id").Where("peoples.name LIKE ?", "%"+*qry.OwnerName+"%")
	}
	if qry.OwnerSurname != nil {
		builder.Joins("JOIN peoples ON peoples.id = cars.owner_id").Where("peoples.surname = ?", "%"+*qry.OwnerSurname+"%")
	}
	builder.Preload("Owner")
	builder = builder.Order(fmt.Sprintf("id %s", qry.Order))
	builder = builder.Limit(qry.Count).Offset((qry.Page - 1) * qry.Count)
	var entities []Car
	result := builder.Find(&entities)
	if result.Error != nil {
		log.Error("failed to search cars", slog.String("error", result.Error.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrDatabase, result.Error)
	}
	cars := make([]model.Car, len(entities))
	for i, entity := range entities {
		cars[i] = ToModel(entity)
	}

	log.Debug("searched cars", slog.Any("cars", cars))

	return &cars, nil
}

func (r *Repository) Create(qry *query.CarCreate) (*model.Car, error) {
	const op = "repository.gorm.car.Create"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("qry", qry),
	)

	log.Info("creating car")

	entity := Car{
		RegNum:  qry.RegNum,
		Mark:    qry.Mark,
		Model:   qry.Model,
		OwnerID: qry.OwnerID,
	}
	if qry.Year != nil {
		entity.Year = *qry.Year
	}
	result := r.db.Create(&entity)
	if result.Error != nil {
		log.Error("failed to create car", slog.String("error", result.Error.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrDatabase, result.Error)
	}

	var fullEntity Car
	if err := r.db.Preload("Owner").First(&fullEntity, entity.ID).Error; err != nil {
		log.Error("failed to load car with owner", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrDatabase, err)
	}
	car := ToModel(fullEntity)

	log.Debug("created car", slog.Any("car", car))

	return &car, nil
}

func (r *Repository) Update(qry *query.CarUpdate) (*model.Car, error) {
	const op = "repository.gorm.car.Update"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("qry", qry),
	)

	log.Info("searching car")

	var entity Car
	result := r.db.Preload("Owner").First(&entity, qry.ID)
	if result.Error != nil {
		log.Error("failed to search", slog.String("error", result.Error.Error()))
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s - %d", app_error.ErrNotFound, "failed to search by id", qry.ID)
		}
		return nil, fmt.Errorf("%w: %w", app_error.ErrDatabase, result.Error)
	}

	log.Debug("searched car", slog.Any("entity", entity))
	log.Info("updating car", slog.Any("entity", entity))

	if qry.RegNum != nil {
		entity.RegNum = *qry.RegNum
	}
	if qry.Mark != nil {
		entity.Mark = *qry.Mark
	}
	if qry.Model != nil {
		entity.Model = *qry.Model
	}
	if qry.Year != nil {
		entity.Year = *qry.Year
	}
	result = r.db.Save(&entity)
	if result.Error != nil {
		log.Error("failed to update car", slog.String("error", result.Error.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrDatabase, result.Error)
	}
	car := ToModel(entity)

	log.Debug("updated car", slog.Any("car", car))

	return &car, nil
}

func (r *Repository) Delete(qry *query.CarDelete) error {
	const op = "repository.gorm.car.Delete"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("qry", qry),
	)

	log.Info("deleting car")

	result := r.db.Delete(&Car{}, qry.ID)
	if result.RowsAffected == 0 {
		log.Error("failed to delete car")
		return fmt.Errorf("%w: %s - %d", app_error.ErrNotFound, "failed to delete by id", qry.ID)
	}
	if result.Error != nil {
		log.Error("failed to delete car", slog.String("error", result.Error.Error()))
		return fmt.Errorf("%w: %w", app_error.ErrDatabase, result.Error)
	}

	log.Debug("deleted car")

	return nil
}

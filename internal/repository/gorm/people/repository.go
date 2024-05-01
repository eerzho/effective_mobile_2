package people

import (
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

func (r *Repository) Create(qry *query.PeopleCreate) (*model.People, error) {
	const op = "repository.gorm.people.Create"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("qry", qry),
	)

	log.Info("creating people")

	entity := People{
		Name:    qry.Name,
		Surname: qry.Surname,
	}
	if qry.Patronymic != nil {
		entity.Patronymic = *qry.Patronymic
	}
	result := r.db.Create(&entity)
	if result.Error != nil {
		log.Error("failed to create people", slog.String("error", result.Error.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrDatabase, result.Error)
	}
	people := ToModel(entity)

	log.Debug("created people", slog.Any("people", people))

	return &people, nil
}

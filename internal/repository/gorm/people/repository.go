package people

import (
	"log/slog"

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
	const op = "repository.people.Create"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("repository starting", slog.Any("qry", qry))

	log.Info("create people")
	people := People{
		Name:    qry.Name,
		Surname: qry.Surname,
	}
	if qry.Patronymic != nil {
		people.Patronymic = *qry.Patronymic
	}
	result := r.db.Create(&people)
	if result.Error != nil {
		log.Error("failed to create", slog.String("error", result.Error.Error()))
		return nil, result.Error
	}
	log.Debug("created people", slog.Any("people", people))

	log.Info("adapting entity to model")
	m := ToModel(people)
	log.Debug("adapted entity to model", slog.Any("people", people))

	return &m, nil
}

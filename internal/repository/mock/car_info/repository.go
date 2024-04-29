package car_info

import (
	"log/slog"
	"time"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/query"
)

type Repository struct {
}

func New() *Repository {
	return &Repository{}
}

func (r *Repository) GetCarInfo(qry *query.CarInfo) (*model.CarInfo, error) {
	const op = "repository.carInfo.GetCarInfo"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("repository starting", slog.Any("qry", qry))

	t := time.Now()
	log.Info("creating mock owner")
	owner := model.People{
		Name:    "Test Owner " + t.String(),
		Surname: "Super",
	}
	log.Debug("created mock owner", slog.Any("owner", owner))

	log.Info("creating mock car info")
	year := t.Year()
	carInfo := model.CarInfo{
		Mark:  "Test Mark " + t.String(),
		Model: "Super",
		Owner: &owner,
		Year:  &year,
	}
	log.Debug("created mock car info", slog.Any("carInfo", carInfo))

	return &carInfo, nil
}

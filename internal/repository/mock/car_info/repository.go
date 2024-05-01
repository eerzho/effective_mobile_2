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
	const op = "repository.mock.car_info.GetCarInfo"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("qry", qry),
	)

	log.Info("getting car info")

	t := time.Now()
	owner := model.People{
		Name:    "Test Owner " + t.String(),
		Surname: "Super",
	}
	year := t.Year()
	carInfo := model.CarInfo{
		Mark:  "Test Mark " + t.String(),
		Model: "Super",
		Owner: &owner,
		Year:  &year,
	}

	log.Debug("got car info", slog.Any("carInfo", carInfo))

	return &carInfo, nil
}

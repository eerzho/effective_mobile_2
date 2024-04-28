package car

import (
	"log/slog"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/command"
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/query"
)

type Repository interface {
	List(qry *query.CarList) (*[]model.Car, error)
	Update(qry *query.CarUpdate) (*model.Car, error)
	Delete(qry *query.CarDelete) error
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Index(cmd *command.CarIndex) (*[]model.Car, error) {
	const op = "service.car.Index"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("service starting", slog.Any("cmd", cmd))

	log.Info("parsing service command")
	qry := query.CarList{
		RegNum: cmd.RegNum,
		Mark:   cmd.Mark,
		Model:  cmd.Model,
		Year:   cmd.Year,
	}
	if cmd.Page == nil || *cmd.Page <= 0 {
		qry.Page = 1
	} else {
		qry.Page = *cmd.Page
	}
	if cmd.Count == nil || *cmd.Count <= 0 {
		qry.Count = 10
	} else {
		qry.Count = *cmd.Count
	}
	log.Debug("parsed command", slog.Any("qry", qry))

	log.Info("executing repository")
	cars, err := s.repository.List(&qry)
	log.Debug("repository result",
		slog.Any("cars", *cars),
		slog.Any("err", err),
	)
	if err != nil {
		log.Error("repository error", slog.String("error", err.Error()))
		return nil, err
	}

	return cars, nil
}

func (s *Service) Update(cmd *command.CarUpdate) (*model.Car, error) {
	const op = "service.car.Update"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("service starting", slog.Any("cmd", cmd))

	log.Info("parsing service command")
	qry := query.CarUpdate{
		ID:     cmd.ID,
		RegNum: cmd.RegNum,
		Mark:   cmd.Mark,
		Model:  cmd.Model,
		Year:   cmd.Year,
	}
	log.Debug("parsed command", slog.Any("qry", qry))

	log.Info("executing repository")
	car, err := s.repository.Update(&qry)
	log.Debug("repository result",
		slog.Any("car", car),
		slog.Any("err", err),
	)
	if err != nil {
		log.Error("repository error", slog.String("error", err.Error()))
		return nil, err
	}

	return car, nil
}

func (s *Service) Delete(cmd *command.CarDelete) error {
	const op = "service.Car.Delete"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("service starting", slog.Any("cmd", cmd))

	log.Info("parsing service command")
	qry := query.CarDelete{ID: cmd.ID}
	log.Debug("parsed command", slog.Any("qry", qry))

	log.Info("executing repository")
	err := s.repository.Delete(&qry)
	log.Debug("repository result", slog.Any("err", err))
	if err != nil {
		log.Error("repository error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

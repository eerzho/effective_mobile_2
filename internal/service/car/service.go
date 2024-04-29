package car

import (
	"log/slog"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/command"
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/query"
)

type Service struct {
	carRepository     carRepository
	carInfoRepository carInfoRepository
	ownerRepository   ownerRepository
}

func New(
	carRepository carRepository,
	carInfoRepository carInfoRepository,
	ownerRepository ownerRepository,
) *Service {
	return &Service{
		carRepository:     carRepository,
		carInfoRepository: carInfoRepository,
		ownerRepository:   ownerRepository,
	}
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
	cars, err := s.carRepository.List(&qry)
	log.Debug("repository result",
		slog.Any("cars", cars),
		slog.Any("err", err),
	)
	if err != nil {
		log.Error("repository error", slog.String("error", err.Error()))
		return nil, err
	}

	return cars, nil
}

func (s *Service) Store(cmd *command.CarStore) (*[]model.Car, error) {
	const op = "service.car.Store"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("service starting", slog.Any("cmd", cmd))

	cars := make([]model.Car, len(cmd.RegNums))
	for i, regNum := range cmd.RegNums {
		log.Info("parsing service command")
		qryCarInfo := query.CarInfo{RegNum: regNum}
		log.Info("executing repository: getting car info")
		carInfo, err := s.carInfoRepository.GetCarInfo(&qryCarInfo)
		log.Debug("repository result",
			slog.Any("carInfo", carInfo),
			slog.Any("err", err),
		)
		if err != nil {
			log.Error("repository error", slog.String("error", err.Error()))
			return nil, err
		}

		log.Info("parsing service command")
		qryPeopleCreate := query.PeopleCreate{
			Name:       carInfo.Owner.Name,
			Surname:    carInfo.Owner.Surname,
			Patronymic: carInfo.Owner.Patronymic,
		}
		log.Info("executing repository: creating people")
		people, err := s.ownerRepository.Create(&qryPeopleCreate)
		log.Debug("repository result",
			slog.Any("people", people),
			slog.Any("err", err),
		)
		if err != nil {
			log.Error("repository error", slog.String("error", err.Error()))
			return nil, err
		}

		log.Info("parsing service command")
		qryCarCreate := query.CarCreate{
			RegNum:  regNum,
			Mark:    carInfo.Mark,
			Model:   carInfo.Model,
			Year:    carInfo.Year,
			OwnerID: people.ID,
		}
		log.Info("executing repository: creating car")
		car, err := s.carRepository.Create(&qryCarCreate)
		log.Debug("repository result",
			slog.Any("car", car),
			slog.Any("err", err),
		)
		if err != nil {
			log.Error("repository error", slog.String("error", err.Error()))
			return nil, err
		}

		cars[i] = *car
	}

	return &cars, nil
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
	car, err := s.carRepository.Update(&qry)
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
	err := s.carRepository.Delete(&qry)
	log.Debug("repository result", slog.Any("err", err))
	if err != nil {
		log.Error("repository error", slog.String("error", err.Error()))
		return err
	}

	return nil
}

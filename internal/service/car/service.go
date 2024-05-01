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
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("cmd", cmd),
	)

	log.Info("searching cars")

	qry := query.CarList{
		RegNum:       cmd.RegNum,
		Mark:         cmd.Mark,
		Model:        cmd.Model,
		Year:         cmd.Year,
		OwnerName:    cmd.OwnerName,
		OwnerSurname: cmd.OwnerSurname,
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
	if cmd.Order == nil || *cmd.Order == "" || (*cmd.Order != "asc" && *cmd.Order != "desc") {
		qry.Order = "desc"
	} else {
		qry.Order = *cmd.Order
	}
	cars, err := s.carRepository.List(&qry)
	if err != nil {
		log.Error("failed to search cars", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("searched cars", slog.Any("cars", cars))

	return cars, nil
}

func (s *Service) Store(cmd *command.CarStore) (*[]model.Car, error) {
	const op = "service.car.Store"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("cmd", cmd),
	)

	log.Info("creating cars")

	cars := make([]model.Car, len(cmd.RegNums))
	for i, regNum := range cmd.RegNums {
		qryCarInfo := query.CarInfo{RegNum: regNum}
		carInfo, err := s.carInfoRepository.GetCarInfo(&qryCarInfo)
		if err != nil {
			log.Error("failed to get car info", slog.String("error", err.Error()))
			return nil, err
		}
		qryPeopleCreate := query.PeopleCreate{
			Name:       carInfo.Owner.Name,
			Surname:    carInfo.Owner.Surname,
			Patronymic: carInfo.Owner.Patronymic,
		}
		people, err := s.ownerRepository.Create(&qryPeopleCreate)
		if err != nil {
			log.Error("failed to create people", slog.String("error", err.Error()))
			return nil, err
		}
		qryCarCreate := query.CarCreate{
			RegNum:  regNum,
			Mark:    carInfo.Mark,
			Model:   carInfo.Model,
			Year:    carInfo.Year,
			OwnerID: people.ID,
		}
		car, err := s.carRepository.Create(&qryCarCreate)
		if err != nil {
			log.Error("failed to create car", slog.String("error", err.Error()))
			return nil, err
		}
		cars[i] = *car
	}

	log.Debug("created cars", slog.Any("cars", cars))

	return &cars, nil
}

func (s *Service) Update(cmd *command.CarUpdate) (*model.Car, error) {
	const op = "service.car.Update"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("cmd", cmd),
	)

	log.Info("updating car")

	qry := query.CarUpdate{
		ID:     cmd.ID,
		RegNum: cmd.RegNum,
		Mark:   cmd.Mark,
		Model:  cmd.Model,
		Year:   cmd.Year,
	}
	car, err := s.carRepository.Update(&qry)
	if err != nil {
		log.Error("failed to update car", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("updated car", slog.Any("car", car))

	return car, nil
}

func (s *Service) Delete(cmd *command.CarDelete) error {
	const op = "service.car.Delete"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("cmd", cmd),
	)

	log.Info("deleting car")

	qry := query.CarDelete{ID: cmd.ID}
	err := s.carRepository.Delete(&qry)
	if err != nil {
		log.Error("failed to delete car", slog.String("error", err.Error()))
		return err
	}

	log.Debug("deleted car")

	return nil
}

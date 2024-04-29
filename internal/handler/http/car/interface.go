package car

import (
	"effective_mobile_2/internal/dto/command"
	"effective_mobile_2/internal/dto/model"
)

type service interface {
	Index(cmd *command.CarIndex) (*[]model.Car, error)
	Store(cmd *command.CarStore) (*[]model.Car, error)
	Update(cmd *command.CarUpdate) (*model.Car, error)
	Delete(cmd *command.CarDelete) error
}

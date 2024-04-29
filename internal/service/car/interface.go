package car

import (
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/query"
)

type carRepository interface {
	List(qry *query.CarList) (*[]model.Car, error)
	Create(qry *query.CarCreate) (*model.Car, error)
	Update(qry *query.CarUpdate) (*model.Car, error)
	Delete(qry *query.CarDelete) error
}

type carInfoRepository interface {
	GetCarInfo(qry *query.CarInfo) (*model.CarInfo, error)
}

type ownerRepository interface {
	Create(qry *query.PeopleCreate) (*model.People, error)
}

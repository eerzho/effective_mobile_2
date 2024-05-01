package car

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/command"
	"effective_mobile_2/internal/handler/http/dto/request"
	"effective_mobile_2/internal/handler/http/dto/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type Handler struct {
	service service
}

func New(service service) *Handler {
	return &Handler{service: service}
}

// Index lists all cars based on query parameters
// @Summary List all cars
// @Description Get a list of cars filtered by various parameters
// @Tags cars
// @Accept json
// @Produce json
// @Param regNum query string false "Registration Number filter"
// @Param mark query string false "Car mark filter"
// @Param model query string false "Car model filter"
// @Param year query int false "Car year filter"
// @Param ownerName query string false "Owner name filter"
// @Param ownerSurname query string false "Owner surname filter"
// @Param order query string false "Order of results (asc or desc)"
// @Param page query int false "Page number for pagination"
// @Param count query int false "Number of items per page"
// @Success 200 {array} model.Car
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/cars [get]
func (h *Handler) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Index"
		log := app_log.Logger().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("searching cars")

		var req request.CarIndex
		if err := schema.NewDecoder().Decode(&req, r.URL.Query()); err != nil {
			log.Error("failed to decode", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}

		cmd := command.CarIndex{
			RegNum:       req.RegNum,
			Mark:         req.Mark,
			Model:        req.Model,
			Year:         req.Year,
			OwnerName:    req.OwnerName,
			OwnerSurname: req.OwnerSurname,
			Order:        req.Order,
			Page:         req.Page,
			Count:        req.Count,
		}
		cars, err := h.service.Index(&cmd)
		if err != nil {
			log.Error("failed to search cars", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}

		log.Debug("searched cars", slog.Any("cars", cars))

		response.Ok(&w, r, cars)
	}
}

// Store creates new cars based on the provided data
// @Summary Create new cars
// @Description Add one or more new cars to the database
// @Tags cars
// @Accept json
// @Produce json
// @Param request body request.CarStore true "New car details"
// @Success 200 {array} model.Car
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/cars [post]
func (h *Handler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Store"
		log := app_log.Logger().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("creating cars")

		var req request.CarStore
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}
		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}

		cmd := command.CarStore{RegNums: req.RegNums}
		cars, err := h.service.Store(&cmd)
		if err != nil {
			log.Error("failed to create cars", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}

		log.Debug("created cars", slog.Any("cars", cars))

		response.Ok(&w, r, cars)
	}
}

// Update modifies an existing car
// @Summary Update car details
// @Description Update details of an existing car by its ID
// @Tags cars
// @Accept json
// @Produce json
// @Param id path int true "Car ID"
// @Param request body request.CarUpdate true "Car update details"
// @Success 200 {object} model.Car
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/cars/{id} [patch]
func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Update"
		log := app_log.Logger().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("updating car")

		var req request.CarUpdate
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}
		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("failed to convert", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}

		cmd := command.CarUpdate{
			ID:     id,
			RegNum: req.RegNum,
			Mark:   req.Mark,
			Model:  req.Model,
			Year:   req.Year,
		}
		car, err := h.service.Update(&cmd)
		if err != nil {
			log.Error("failed to update car", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}

		log.Debug("updated car", slog.Any("car", car))

		response.Ok(&w, r, car)
	}
}

// Delete removes a car
// @Summary Remove a car
// @Description Delete a car by its ID
// @Tags cars
// @Accept json
// @Produce json
// @Param id path int true "Car ID"
// @Success 200
// @Failure 400 {object} response.Error
// @Failure 404 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/cars/{id} [delete]
func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Delete"
		log := app_log.Logger().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("deleting car")

		var cmd command.CarDelete
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("failed to convert", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}
		cmd.ID = id
		if err = h.service.Delete(&cmd); err != nil {
			log.Error("failed to delete car", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}

		log.Debug("deleted car")

		response.Ok(&w, r, nil)
	}
}

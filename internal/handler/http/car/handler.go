package car

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/command"
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

func (h *Handler) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Index"
		log := app_log.Logger().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("searching cars")

		var cmd command.CarIndex
		if err := schema.NewDecoder().Decode(&cmd, r.URL.Query()); err != nil {
			log.Error("failed to decode", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
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

func (h *Handler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Store"
		log := app_log.Logger().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("creating cars")

		var cmd command.CarStore
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Error("failed to decode", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}
		if err := validator.New().Struct(cmd); err != nil {
			log.Error("failed to validate", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}
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

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Update"
		log := app_log.Logger().With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("updating car")

		var cmd command.CarUpdate
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("failed to convert", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}
		cmd.ID = id
		if err = json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Error("failed to decode", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
		}
		if err = validator.New().Struct(cmd); err != nil {
			log.Error("failed to validate", slog.String("error", err.Error()))
			response.Bad(&w, r, err)
			return
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

		response.Ok(&w, r, "")
	}
}

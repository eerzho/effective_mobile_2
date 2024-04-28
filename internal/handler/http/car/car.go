package car

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/command"
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type Service interface {
	Index(cmd *command.CarIndex) (*[]model.Car, error)
	Update(cmd *command.CarUpdate) (*model.Car, error)
	Delete(cmd *command.CarDelete) error
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Index"
		log := app_log.Logger().With(slog.String("op", op))

		log.Debug("starting handler", slog.String("request_id", middleware.GetReqID(r.Context())))

		log.Info("parsing query")
		var cmd command.CarIndex
		if err := schema.NewDecoder().Decode(&cmd, r.URL.Query()); err != nil {
			log.Error("failed to parse query", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error{Error: err.Error()})
			return
		}
		log.Debug("parsed query", slog.Any("cmd", cmd))

		log.Info("executing service")
		cars, err := h.service.Index(&cmd)
		log.Debug("service result",
			slog.Any("cars", cars),
			slog.Any("err", err),
		)
		if err != nil {
			log.Error("service error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Success{Data: cars})
	}
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Update"
		log := app_log.Logger().With(slog.String("op", op))

		log.Debug("starting handler", slog.String("request_id", middleware.GetReqID(r.Context())))

		log.Info("parsing url")
		var cmd command.CarUpdate
		idParam := chi.URLParam(r, "id")
		if id, err := strconv.Atoi(idParam); err != nil {
			log.Error("failed to parse url", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error{Error: err.Error()})
			return
		} else {
			cmd.ID = id
		}
		log.Debug("parsed url", slog.Any("cmd", cmd))

		log.Info("parsing body")
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Error("failed to parse body", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			if errors.Is(err, io.EOF) {
				render.JSON(w, r, response.Error{Error: "empty body"})
				return
			}
			render.JSON(w, r, response.Error{Error: err.Error()})
			return
		}
		log.Debug("parsed body", slog.Any("cmd", cmd))

		log.Info("validating body")
		if err := validator.New().Struct(cmd); err != nil {
			log.Error("failed to validate body", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error{Error: err.Error()})
			return
		}
		log.Debug("validated body", slog.Any("cmd", cmd))

		log.Info("executing service")
		car, err := h.service.Update(&cmd)
		log.Debug("service result",
			slog.Any("car", car),
			slog.Any("err", err),
		)
		if err != nil {
			log.Error("service error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Success{Data: car})
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.http.car.Delete"
		log := app_log.Logger().With(slog.String("op", op))

		log.Debug("starting handler", slog.String("request_id", middleware.GetReqID(r.Context())))

		log.Info("parsing url")
		var cmd command.CarDelete
		idParam := chi.URLParam(r, "id")
		if id, err := strconv.Atoi(idParam); err != nil {
			log.Error("failed to parse url", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error{Error: err.Error()})
			return
		} else {
			cmd.ID = id
		}
		log.Debug("parsed url", slog.Any("cmd", cmd))

		log.Info("executing service")
		err := h.service.Delete(&cmd)
		log.Debug("service result", slog.Any("err", err))
		if err != nil {
			log.Error("service error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Success{})
	}
}

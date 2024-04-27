package car

import (
	"log/slog"
	"net/http"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/command"
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

type Service interface {
	Index(cmd command.CarIndex) (*[]model.Car, error)
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

		log.Info("parsing request query")
		var cmd command.CarIndex
		if err := schema.NewDecoder().Decode(&cmd, r.URL.Query()); err != nil {
			log.Error("failed to parse query", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error{Error: err.Error()})
		}
		log.Debug("parsed query", slog.Any("cmd", cmd))

		log.Info("executing service")
		cars, err := h.service.Index(cmd)
		log.Debug("service result",
			slog.Any("cars", cars),
			slog.Any("err", err),
		)
		if err != nil {
			log.Error("service error", slog.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error{Error: err.Error()})
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response.Success{Data: cars})
	}
}

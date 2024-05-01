package http_srv

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "effective_mobile_2/docs"
	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/config"
	"effective_mobile_2/internal/database"
	carH "effective_mobile_2/internal/handler/http/car"
	carGR "effective_mobile_2/internal/repository/gorm/car"
	peopleGR "effective_mobile_2/internal/repository/gorm/people"
	httpSwagger "github.com/swaggo/http-swagger"

	carInfoAR "effective_mobile_2/internal/repository/api/car_info"
	//carInfoMock "effective_mobile_2/internal/repository/mock/car_info"
	carS "effective_mobile_2/internal/service/car"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() error {
	const op = "app.http.Run"

	log := app_log.Logger().With(slog.String("op", op))

	log.Info("configuring http server")
	router := chi.NewRouter()

	setupMiddleware(router)
	setupEndpoints(router)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	httpSrv := http.Server{
		Addr:    config.Cfg().Http.Address,
		Handler: router,
	}

	log.Info("serving http server")
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			app_log.Logger().Error("failed to start http server", slog.String("error", err.Error()))
		}
	}()

	log.Info("http server listening on " + config.Cfg().Http.Address)

	<-done

	log.Info("shutting down http server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func setupMiddleware(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use()
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func setupEndpoints(router *chi.Mux) {

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), // The url pointing to API definition
	))

	carRepository := carGR.New(database.Db().Gorm)
	carInfoRepository := carInfoAR.New(config.Cfg().Api.CarInfo)
	//carInfoRepository := carInfoMock.New()
	peopleRepository := peopleGR.New(database.Db().Gorm)

	carService := carS.New(carRepository, carInfoRepository, peopleRepository)

	carHandler := carH.New(carService)

	router.Get("/api/cars", carHandler.Index())
	router.Post("/api/cars", carHandler.Store())
	router.Patch("/api/cars/{id}", carHandler.Update())
	router.Delete("/api/cars/{id}", carHandler.Delete())
}

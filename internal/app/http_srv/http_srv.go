package http_srv

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/config"
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
		Addr:         config.Cfg().Http.Address,
		Handler:      router,
		ReadTimeout:  config.Cfg().Http.Timeout,
		WriteTimeout: config.Cfg().Http.Timeout,
		IdleTimeout:  config.Cfg().Http.IdleTimeout,
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

}

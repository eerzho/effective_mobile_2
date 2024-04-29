package car_info

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"effective_mobile_2/internal/app_log"
	"effective_mobile_2/internal/dto/model"
	"effective_mobile_2/internal/dto/query"
)

type Repository struct {
	url string
}

func New(url string) *Repository {
	return &Repository{url: strings.TrimRight(url, "/")}
}

func (r *Repository) GetCarInfo(qry *query.CarInfo) (*model.CarInfo, error) {
	const op = "repository.carInfo.GetCarInfo"
	log := app_log.Logger().With(slog.String("op", op))

	log.Debug("repository starting", slog.Any("qry", qry))

	log.Info("getting car info")
	url := r.url + "/info?regNum=" + qry.RegNum

	log.Info("creating request")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("failed to create http request", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("created request", slog.Any("req", req))

	log.Info("sending request")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed to send http request", slog.String("error", err.Error()))
		return nil, err
	}
	defer resp.Body.Close()
	log.Debug("sent request", slog.Any("resp", resp))

	log.Info("validating response")
	if resp.StatusCode != http.StatusOK {
		log.Error("failed to validate response", slog.String("status", resp.Status))
		return nil, fmt.Errorf("failed to validate response: %s", resp.Status)
	}
	log.Debug("validated response", slog.Int("statusCode", resp.StatusCode))

	log.Info("reading response")
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("read response", slog.Any("responseBody", responseBody))

	log.Info("parsing response")
	carInfo := &model.CarInfo{}
	if err := json.Unmarshal(responseBody, &carInfo); err != nil {
		log.Error("failed to parse response", slog.String("error", err.Error()))
		return nil, err
	}
	log.Debug("parsed response", slog.Any("carInfo", carInfo))

	return carInfo, nil
}

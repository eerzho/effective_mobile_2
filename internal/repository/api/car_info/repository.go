package car_info

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"effective_mobile_2/internal/app_error"
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
	const op = "repository.api.car_info.GetCarInfo"
	log := app_log.Logger().With(
		slog.String("op", op),
		slog.Any("qry", qry),
	)

	log.Info("getting car info")

	req, err := http.NewRequest("GET", r.url+"/info?regNum="+qry.RegNum, nil)
	if err != nil {
		log.Error("failed to create HTTP request", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrHTTPRequestFailed, err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("HTTP request failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrHTTPRequestFailed, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("HTTP request failed with status", slog.Int("status", resp.StatusCode))
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("%w: %s - %s", app_error.ErrNotFound, "failed to get car info by regNum", qry.RegNum)
		}
		return nil, fmt.Errorf("%w: %s - %d", app_error.ErrHTTPRequestFailed, "request failed with status", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrHTTPRequestFailed, err)
	}

	carInfo := model.CarInfo{}
	if err = json.Unmarshal(responseBody, &carInfo); err != nil {
		log.Error("JSON unmarshaling failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %w", app_error.ErrHTTPRequestFailed, err)
	}

	log.Debug("got car info", slog.Any("carInfo", carInfo))

	return &carInfo, nil
}

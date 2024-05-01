package response

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"effective_mobile_2/internal/app_error"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Error struct {
	Message string `json:"message"`
}

func Bad(w *http.ResponseWriter, r *http.Request, err error) {
	var code int
	var message string

	switch {
	case errors.Is(err, app_error.ErrNotFound):
		code = http.StatusNotFound
		message = err.Error()
	case errors.As(err, &validator.ValidationErrors{}):
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			field := ve[0].StructField()
			tag := ve[0].Tag()
			message = fmt.Sprintf("'%s': must be %s", field, tag)
			code = http.StatusBadRequest
		}
	case err == io.EOF:
		code = http.StatusBadRequest
		message = "body is empty"
	default:
		code = http.StatusInternalServerError
		message = err.Error()
	}

	(*w).WriteHeader(code)
	render.JSON(*w, r, Error{Message: message})
}

func Ok(w *http.ResponseWriter, r *http.Request, data interface{}) {
	(*w).WriteHeader(http.StatusOK)
	render.JSON(*w, r, data)
}

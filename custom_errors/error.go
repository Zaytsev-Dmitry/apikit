package custom_errors

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	RowNotFound     = errors.New("запись не найдена")
	ValidationError = errors.New("ошибка валидации")
	MarshallError   = errors.New("ошибка маршалинга")
	ConflictError   = errors.New("запись с такими данными уже существует")
	ForbiddenError  = errors.New("доступ запрещен")
	Unauthorized    = errors.New("нет прав доступа")
)

type MetaData struct {
	Path      string `json:"path"`
	Timestamp string `json:"timestamp"`
}

type ErrorBuilderFunc[T any] func(description string, code int, meta MetaData) T

func HandleError[T any](c *gin.Context, err error, builder ErrorBuilderFunc[T]) {
	resp, status := buildErrorResponse(c, err, builder)
	c.JSON(status, resp)
}

func SetResponseError[T any](c *gin.Context, err error, builder ErrorBuilderFunc[T]) {
	status, msg := getErrorMsgAndStatus(err)
	resp := builder(msg, status, getMeta(c))

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Status(status)
	json.NewEncoder(c.Writer).Encode(resp)
}

func getMeta(c *gin.Context) MetaData {
	return MetaData{
		Path:      c.Request.URL.Path,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func getErrorMsgAndStatus(err error) (int, string) {
	switch {
	case errors.Is(err, RowNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, MarshallError), errors.Is(err, ValidationError):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, ConflictError):
		return http.StatusConflict, err.Error()
	case errors.Is(err, Unauthorized):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, ForbiddenError):
		return http.StatusForbidden, err.Error()
	default:
		return http.StatusInternalServerError, "Oops... что-то пошло не так"
	}
}

func buildErrorResponse[T any](c *gin.Context, err error, builder ErrorBuilderFunc[T]) (T, int) {
	status, msg := getErrorMsgAndStatus(err)
	meta := getMeta(c)
	return builder(msg, status, meta), status
}

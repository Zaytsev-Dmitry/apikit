package custom_errors

import (
	"encoding/json"
	"errors"
	"github.com/Zaytsev-Dmitry/apikit/dto"
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

func HandleError(c *gin.Context, err error) {
	status, msg := getErrorMsgAndStatus(err)
	responseError := getErrorDto(msg, status, c)
	c.JSON(status, responseError)
}

func GetErrorDto(c *gin.Context, err error) dto.BackendErrorResponse {
	status, msg := getErrorMsgAndStatus(err)
	return getErrorDto(msg, status, c)
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

func getErrorDto(err string, errorCode int, context *gin.Context) dto.BackendErrorResponse {
	nowString := time.Now().String()
	return dto.BackendErrorResponse{
		Description: &err,
		ErrorCode:   &errorCode,
		Meta: &dto.MetaData{
			Path:      &context.Request.URL.Path,
			Timestamp: &nowString,
		},
	}
}

func SetResponseError(context *gin.Context, err error) {
	var status int
	msg := "Oops...что то пошло не так"

	switch {
	case errors.Is(err, RowNotFound):
		status = http.StatusNotFound
		msg = err.Error()
	case errors.Is(err, MarshallError) || errors.Is(err, ValidationError):
		status = http.StatusBadRequest
		msg = err.Error()
	default:
		status = http.StatusInternalServerError
	}

	var responseError = getErrorDto(msg, status, context)
	context.Writer.Header().Set("Content-Type", "application/json")
	context.Status(status)
	json.NewEncoder(context.Writer).Encode(responseError)
}

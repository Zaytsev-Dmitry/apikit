package handlers

import (
	"github.com/Zaytsev-Dmitry/apikit/custom_errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleMarshalling[T any](context *gin.Context, req *T) error {
	if err := context.ShouldBindJSON(req); err != nil {
		custom_errors.HandleError(context, custom_errors.MarshallError)
	}
	return nil
}

func HandleResponse[T any, R any](c *gin.Context, logic func() (T, error), present func(T, *gin.Context) R) {
	if result, err := logic(); err != nil {
		custom_errors.HandleError(c, err)
	} else {
		c.JSON(http.StatusOK, present(result, c))
	}
}

func HandleResponseWithoutPresent(c *gin.Context, logic func() error) {
	if err := logic(); err != nil {
		custom_errors.HandleError(c, err)
	} else {
		c.Status(http.StatusNoContent)
	}
}

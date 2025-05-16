package handlers

import (
	"github.com/Zaytsev-Dmitry/apikit/custom_errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleMarshalling[T any](c *gin.Context, req *T) error {
	if err := c.ShouldBindJSON(req); err != nil {
		custom_errors.SetResponseError(c, custom_errors.MarshallError)
		return err
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

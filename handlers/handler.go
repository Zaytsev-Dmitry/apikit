package handlers

import (
	"github.com/Zaytsev-Dmitry/apikit/custom_errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleMarshalling[T any](c *gin.Context, req *T, errorBuilder custom_errors.ErrorBuilderFunc[any]) error {
	if err := c.ShouldBindJSON(req); err != nil {
		custom_errors.SetResponseError(c, custom_errors.MarshallError, errorBuilder)
		return err
	}
	return nil
}

func HandleResponse[T any, R any, E any](
	c *gin.Context,
	logic func() (T, error),
	present func(T, *gin.Context) R,
	errorBuilder custom_errors.ErrorBuilderFunc[E],
) {
	if result, err := logic(); err != nil {
		custom_errors.HandleError(c, err, errorBuilder)
	} else {
		c.JSON(http.StatusOK, present(result, c))
	}
}

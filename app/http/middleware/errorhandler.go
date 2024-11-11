package middleware

import (
	"errors"
	"github.com/TicketsBot/GoPanel/app"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		var message string

		var apiError *app.ApiError
		if errors.As(c.Errors[0], &apiError) {
			message = apiError.ExternalMessage
		} else {
			message = "An error occurred processing your request"
		}

		c.JSON(-1, ErrorResponse{
			Error: message,
		})
	}

	if c.Writer.Status() >= 500 && len(c.Errors) == 0 {
		c.JSON(-1, ErrorResponse{
			Error: "An internal server error occurred",
		})
	}
}

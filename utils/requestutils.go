package utils

import "github.com/gin-gonic/gin"

func ErrorToResponse(err error) map[string]interface{} {
	return gin.H {
		"success": false,
		"error": err.Error(),
	}
}

var SuccessResponse = gin.H{
	"success": true,
}

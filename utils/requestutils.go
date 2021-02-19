package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ErrorJson(err error) map[string]interface{} {
	return ErrorStr(err.Error())
}

func ErrorStr(err string, format ...interface{}) map[string]interface{} {
	return gin.H {
		"success": false,
		"error": fmt.Sprintf(err, format...),
	}
}

var SuccessResponse = gin.H{
	"success": true,
}

package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ErrorJson(err error) map[string]any {
	return ErrorStr(err.Error())
}

func ErrorStr(err string, format ...any) map[string]any {
	return gin.H{
		"success": false,
		"error":   fmt.Sprintf(err, format...),
	}
}

var SuccessResponse = gin.H{
	"success": true,
}

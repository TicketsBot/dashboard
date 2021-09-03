package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

func ReadResponse(ctx *gin.Context) {
	ctx.Writer = NewCustomWriter(ctx)
}

type CustomWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func NewCustomWriter(ctx *gin.Context) *CustomWriter {
	return &CustomWriter{
		ResponseWriter: ctx.Writer,
		body:           bytes.NewBuffer([]byte{}),
	}
}

func (w CustomWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *CustomWriter) Read() []byte  {
	return w.body.Bytes()
}

package ginutil

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogWriter struct {
	*zap.SugaredLogger
}

func (l *LogWriter) Write(d []byte) (n int, err error) {
	l.Info(string(d))
	return len(d), nil
}

type NullWriter struct{}

func (l *NullWriter) Write(d []byte) (n int, err error) {
	return len(d), nil
}

func NewGinEngine(out io.Writer, logLevel string) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
		out = os.Stdout
	}

	e.Use(gin.RecoveryWithWriter(out))
	e.Use(gin.LoggerWithWriter(out))

	return e
}

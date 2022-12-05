package log

import (
	"log"
	"testing"

	"go.uber.org/zap"
)

func Test_Logger(t *testing.T) {
	Zap.Info("hello", zap.String("uid", "abcd"))
	err := SetLogLevel("debug")
	if err != nil {
		log.Println(err)
	}
	Zap.Debug("debug log")
}

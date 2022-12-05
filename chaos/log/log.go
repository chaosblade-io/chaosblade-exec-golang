package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap *zap.Logger
var level zap.AtomicLevel
var cfg zapcore.EncoderConfig

func init() {
	cfg := zap.NewProductionEncoderConfig()
	level = zap.NewAtomicLevel()
	Zap = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.Lock(os.Stdout),
		level))
	Zap.Named("chaos")
}

func SetLogLevel(logLevel string) error {
	return level.UnmarshalText([]byte(logLevel))
}

package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func main() {
	lgcf := zap.NewProductionConfig()
	lgcf.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	lgcf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	lgcf.OutputPaths = []string{"/tmp/log.zap", "stderr"}
	logger, err := lgcf.Build()
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i < 30; i++ {
		logger.Debug("debug msg", zap.Int("log cnt", i))
		logger.Info("info msg", zap.Int("log cnt", i))
		logger.Error("error msg", zap.Int("log cnt", i))
		logger.Warn("warn msg", zap.Int("log cnt", i))
		time.Sleep(time.Second)
	}
}

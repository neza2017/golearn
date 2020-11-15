package main

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "/home/cyf/tmp/zap.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)

	//reader := bufio.NewReader(os.Stdin)
	for n := 1; n < 20; n++ {
		for i := 1; i < 200000; i++ {
			logger.Debug("debug msg", zap.Int("log cnt", i), zap.Int("idx", n))
			logger.Info("info msg", zap.Int("log cnt", i), zap.Int("idx", n))
			logger.Error("error msg", zap.Int("log cnt", i), zap.Int("idx", n))
			logger.Warn("warn msg", zap.Int("log cnt", i), zap.Int("idx", n))
		}
		fmt.Printf("=====================================")
		//reader.ReadRune()
	}
}

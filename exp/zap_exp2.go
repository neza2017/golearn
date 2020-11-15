package main

import (
	"bufio"
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type LogWriter struct {
	FileName string
	fp       *os.File
	mu       sync.Mutex
	fileCnt  int
	sig      chan os.Signal
	ctx      context.Context
	cancel   context.CancelFunc
}

func (m *LogWriter) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.fp == nil {
		if err := m.open(); err != nil {
			return 0, nil
		}
	}
	if m.fileCnt == 0 {
		m.fileCnt++
		m.ctx, m.cancel = context.WithCancel(context.Background())
		m.sig = make(chan os.Signal, 1)
		signal.Notify(m.sig, syscall.SIGHUP)
		go func() {
			for {
				select {
				case <-m.ctx.Done():
					return
				case _, ok := <-m.sig:
					if ok == false {
						return
					}
					func() {
						m.mu.Lock()
						defer m.mu.Unlock()
						if m.fp != nil {
							_ = m.fp.Close()
							m.fp = nil
						}
					}()
				}
			}
		}()
	}
	_, _ = os.Stderr.Write(p)
	return m.fp.Write(p)
}

func (m *LogWriter) Sync() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_ = os.Stderr.Sync()
	return m.fp.Sync()
}

func (m *LogWriter) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cancel()
	if m.fp != nil {
		if err := m.fp.Close(); err != nil {
			return err
		}
	}
	if m.sig != nil {
		signal.Reset(syscall.SIGHUP)
		close(m.sig)
	}
	return nil
}

func (m *LogWriter) open() error {
	if m.fp != nil {
		if err := m.fp.Sync(); err != nil {
			return err
		}
		if err := m.fp.Close(); err != nil {
			return err
		}
	}
	fp, err := os.OpenFile(m.FileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	m.fp = fp
	return nil
}

func main() {
	zapid, _ := os.OpenFile("/tmp/zap.pid", os.O_CREATE|os.O_WRONLY, 0644)
	_, _ = zapid.WriteString(fmt.Sprintf("%d", os.Getpid()))
	_ = zapid.Sync()
	_ = syscall.Flock(int(zapid.Fd()), syscall.LOCK_EX)

	mylog := &LogWriter{FileName: "/tmp/zap.log"}
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), mylog, zap.InfoLevel)
	logger := zap.New(core)

	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 10; i++ {
		logger.Debug("debug msg", zap.Int("idx", i))
		logger.Info("info msg", zap.Int("idx", i))
		logger.Error("error msg", zap.Int("idx", i))
		logger.Warn("warn msg", zap.Int("idx", i))
		fmt.Printf("====================")
		_, _, _ = reader.ReadRune()
	}
	_ = mylog.Close()

	_ = syscall.Flock(int(zapid.Fd()), syscall.LOCK_UN)
	_ = zapid.Close()
}

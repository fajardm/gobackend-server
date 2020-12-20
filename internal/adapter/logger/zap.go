package logger

import (
	"fmt"

	"github.com/fajardm/gobackend-server/config"
	"go.uber.org/zap"
)

type Zap struct {
	zap *zap.Logger
}

func NewZap(zap *zap.Logger) *Zap {
	return &Zap{zap}
}

func (l Zap) Error(format string, v ...interface{}) {
	l.zap.Error(fmt.Sprintf(format, v...))
}

func InitZap(conf *config.Config) (*zap.Logger, error) {
	if conf.App.IsDevelopment() {
		logger, err := zap.NewDevelopment(nil)
		if err != nil {
			return nil, err
		}
		return logger, nil
	}
	return nil, nil
}

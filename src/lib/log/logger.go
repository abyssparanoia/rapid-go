package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type loggerKey struct{}

// NewLogger ... loggerを生成し、contextに仕込む。
func NewLogger(ctx context.Context, isDev bool) {

	// 開発環境
	if isDev {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		context.WithValue(ctx, loggerKey{}, logger)
		// production環境
	} else {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		context.WithValue(ctx, loggerKey{}, logger)
	}

}

// Logger ... contextからloggerを取得する
func Logger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(loggerKey{}).(*zap.Logger)

	if logger == nil {
		logger, _ := zap.NewDevelopment()
		return logger
	}

	return logger
}

// Debugf ... Debugログを出力する
func Debugf(ctx context.Context, msg string, fields ...interface{}) {
	Logger(ctx).Debug(fmt.Sprintf(msg, fields...))
}

// Infof ... Infoログを出力する
func Infof(ctx context.Context, msg string, fields ...interface{}) {
	Logger(ctx).Info(fmt.Sprintf(msg, fields...))
}

// Warningf ... Warningログを出力する
func Warningf(ctx context.Context, msg string, fields ...interface{}) {
	Logger(ctx).Warn(fmt.Sprintf(msg, fields...))
}

// Errorf ... Errorログを出力する
func Errorf(ctx context.Context, msg string, fields ...interface{}) {
	Logger(ctx).Error(fmt.Sprintf(msg, fields...))
}

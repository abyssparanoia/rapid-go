package mysql

import (
	"context"

	"github.com/abyssparanoia/rapid-go/src/lib/log"
)

// Logger ... logger for gorm
type Logger struct {
	ctx context.Context
}

// Println ... output sql log
func (logger *Logger) Println(values ...interface{}) {
	texts := ""
	for _, value := range values {
		if text, ok := value.(string); ok {
			texts += text
		}
	}
	log.Infof(logger.ctx, texts)
}

// NewLogger ... get logger for gorm
func NewLogger(ctx context.Context) *Logger {
	return &Logger{ctx}
}

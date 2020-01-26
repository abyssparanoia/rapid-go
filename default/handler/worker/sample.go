package worker

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/internal/pkg/renderer"
)

// SampleHandler ... sample handler
type SampleHandler struct {
}

// Cron ... cron handler
func (h *SampleHandler) Cron(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Debugf(ctx, "call cron handler")
	renderer.Success(ctx, w)
}

// TaskQueue ... task queue handler
func (h *SampleHandler) TaskQueue(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Debugf(ctx, "call task queue handler")
	renderer.Success(ctx, w)
}

// NewSampleHandler ... get sample handler
func NewSampleHandler() *SampleHandler {
	return &SampleHandler{}
}

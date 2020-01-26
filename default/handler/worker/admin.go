package worker

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/internal/pkg/log"
	"github.com/abyssparanoia/rapid-go/internal/pkg/renderer"
)

// AdminHandler ... handler for admin
type AdminHandler struct {
}

// MigrateMasterData ... insert master data
func (h *AdminHandler) MigrateMasterData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Debugf(ctx, "call migrate master data handler")

	renderer.Success(ctx, w)
}

// MigrateTestData ... insert master data
func (h *AdminHandler) MigrateTestData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Debugf(ctx, "call migrate test data handler")

	renderer.Success(ctx, w)
}

// NewAdminHandler ... insert admin handler
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

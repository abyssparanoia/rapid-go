package worker

import (
	"net/http"

	"github.com/abyssparanoia/rapid-go/src/handler"
	"github.com/abyssparanoia/rapid-go/src/lib/log"
)

// AdminHandler ... handler for admin
type AdminHandler struct {
}

// MigrateMasterData ... insert master data
func (h *AdminHandler) MigrateMasterData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Debugf(ctx, "call migrate master data handler")

	handler.RenderSuccess(w)
}

// MigrateTestData ... insert master data
func (h *AdminHandler) MigrateTestData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Debugf(ctx, "call migrate test data handler")

	handler.RenderSuccess(w)
}

// NewAdminHandler ... insert admin handler
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

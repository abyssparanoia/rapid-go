package main

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/environment"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner"
	spanner_repository "github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/transactable"
	"github.com/abyssparanoia/rapid-go/internal/pkg/logger"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/abyssparanoia/rapid-go/internal/pkg/uuid"
	"github.com/caarlos0/env/v10"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/volatiletech/null/v8"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	e := &environment.Environment{}
	if err := env.Parse(e); err != nil {
		panic(err)
	}

	logger := logger.New()
	ctx = ctxzap.ToContext(ctx, logger)

	spannerCli := spanner.NewClient(
		e.SpannerProjectID,
		e.SpannerInstanceID,
		e.SpannerDatabaseID,
	)

	transactable := transactable.NewTransactable(spannerCli)
	tenantRepository := spanner_repository.NewTenant()
	// staffRepository := spanner_repository.NewStaff()

	n := now.Now()
	tenant := model.NewTenant(
		"tenant-name",
		n,
	)
	tenant.ID = uuid.UUID()

	if err := transactable.RWTx(ctx, func(ctx context.Context) error {
		if err := tenantRepository.Create(ctx, tenant); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := transactable.ROTx(ctx, func(ctx context.Context) error {
		got, err := tenantRepository.List(
			ctx,
			repository.ListTenantsQuery{
				BaseListOptions: repository.BaseListOptions{
					Page:  null.Uint64From(1),
					Limit: null.Uint64From(1),
				},
			},
		)
		if err != nil {
			return err
		}
		logger.Info("tenants", zap.Any("tenants", got))
		return nil
	}); err != nil {
		panic(err)
	}
}

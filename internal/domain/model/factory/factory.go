package factory

import (
	"github.com/bxcodec/faker"
	"github.com/playground-live/moala-meet-and-greet-back/internal/domain/model"
	"github.com/playground-live/moala-meet-and-greet-back/internal/pkg/now"
)

func NewFactory() struct {
	Tenant *model.Tenant
} {
	n := now.Now()

	tenant := &model.Tenant{}
	if err := faker.FakeData(tenant); err != nil {
		panic(err)
	}

	tenant.CreatedAt = n
	tenant.UpdatedAt = n

	return struct{ Tenant *model.Tenant }{
		Tenant: tenant,
	}
}

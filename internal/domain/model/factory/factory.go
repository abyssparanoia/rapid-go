package factory

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/bxcodec/faker"
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

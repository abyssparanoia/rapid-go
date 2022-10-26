package factory

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/bxcodec/faker"
)

func NewFactory() struct {
	RequestTime time.Time
	Tenant      *model.Tenant
} {
	n := now.Now()

	tenant := &model.Tenant{}
	if err := faker.FakeData(tenant); err != nil {
		panic(err)
	}

	tenant.CreatedAt = n
	tenant.UpdatedAt = n

	return struct {
		RequestTime time.Time
		Tenant      *model.Tenant
	}{
		RequestTime: n,
		Tenant:      tenant,
	}
}

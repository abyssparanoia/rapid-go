//nolint:exhaustruct
package factory

import (
	"reflect"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/go-faker/faker/v4"
	"github.com/go-faker/faker/v4/pkg/options"
)

func NewFactory() struct {
	RequestTime time.Time
	Tenant      *model.Tenant
	Staff       *model.Staff
} {
	opts := []options.OptionFunc{
		options.WithIgnoreInterface(true),
		options.WithNilIfLenIsZero(true),
		options.WithRandomMapAndSliceMaxSize(1),
	}
	n := now.Now()

	tenant := &model.Tenant{}
	if err := faker.FakeData(tenant, opts...); err != nil {
		panic(err)
	}
	tenant.CreatedAt = n
	tenant.UpdatedAt = n

	user := &model.Staff{}
	if err := faker.FakeData(user, opts...); err != nil {
		panic(err)
	}
	user.TenantID = tenant.ID
	user.Tenant = tenant
	user.CreatedAt = n
	user.UpdatedAt = n

	return struct {
		RequestTime time.Time
		Tenant      *model.Tenant
		Staff       *model.Staff
	}{
		RequestTime: n,
		Tenant:      tenant,
		Staff:       user,
	}
}

func CloneValue(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		x.Interface()
	}
}

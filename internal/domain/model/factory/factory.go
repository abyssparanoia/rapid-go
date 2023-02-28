package factory

import (
	"reflect"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
	"github.com/bxcodec/faker"
)

func NewFactory() struct {
	RequestTime time.Time
	Tenant      *model.Tenant
	Staff       *model.Staff
} {
	n := now.Now()

	tenant := &model.Tenant{}
	if err := faker.FakeData(tenant); err != nil {
		panic(err)
	}
	tenant.CreatedAt = n
	tenant.UpdatedAt = n

	user := &model.Staff{}
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

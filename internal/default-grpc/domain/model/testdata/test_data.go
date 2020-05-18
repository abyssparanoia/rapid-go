package testdata

import (
	"time"

	"github.com/abyssparanoia/rapid-go/internal/default-grpc/domain/model"
	"github.com/bxcodec/faker"
)

// NewDomainModel :
func NewDomainModel() *struct {
	User *model.User
} {

	user := &model.User{}
	if err := faker.FakeData(user); err != nil {
		panic(err)
	}

	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}

	return &struct {
		User *model.User
	}{
		User: user,
	}
}

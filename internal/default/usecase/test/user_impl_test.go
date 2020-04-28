package test

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/default/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/default/domain/model/testdata"
	"github.com/abyssparanoia/rapid-go/internal/default/domain/repository"
	repository_impl "github.com/abyssparanoia/rapid-go/internal/default/infrastructure/repository"
	"github.com/abyssparanoia/rapid-go/internal/default/usecase"
	"github.com/abyssparanoia/rapid-go/internal/testutil"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {

	var ctx context.Context
	var mockCtrl *gomock.Controller
	var userRepository repository.User

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	var _ = Describe("Get", func() {

		var user *model.User

		var prepareData = func() {
			ctx = testutil.Context()
			td := testdata.NewDomainModel()
			user = td.User
		}

		var prepareRepositories = func() {
			userRepository = repository_impl.NewUser()
			_, err := userRepository.Create(ctx, user)
			Expect(err).To(BeNil())
		}

		BeforeEach(func() {
			prepareData()
			prepareRepositories()
		})

		var subject = func() usecase.User {
			return usecase.NewUser(userRepository)
		}

		It("success: get correct user data", func() {
			out, err := subject().Get(ctx, user.ID)
			Expect(err).To(BeNil())
			Expect(out.DisplayName).To(Equal(user.DisplayName))
			Expect(out.Email).To(Equal(user.Email))
		})
	})

})

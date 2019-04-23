
PROJECT_ID = abyssparanoia/rapid-go

init:
	@echo Initialize rapid-go now......
	$(shell go get -u github.com/golang/mock/gomock)
	$(shell go install github.com/golang/mock/mockgen)
	$(shell go get -u github.com/golang/dep/cmd/dep)
	$(shell dep ensure)
	@echo Initialize rapid-go completed!!!!

build:
	$(shell docker-compose build)

start:
	$(shell docker-compose up -d)

down:
	$(shell docker-compose down)

logs:
	$(shell docker-compose logs api)

test:
	$(shell go test -test.v ./src/...)

mockgen_task:
	$(eval SERVICE_LIST := $(call get_service_list))	
	$(foreach file, $(SERVICE_LIST), $(call mockgen_service,$(shell basename $(file))))
	$(eval REPOSITORY_LIST := $(call get_repository_list))
	$(foreach file, $(REPOSITORY_LIST), $(call mockgen_repository,$(shell basename $(file))))

define get_service_list
	$(shell	find ./src/service/ -maxdepth 1 -type f ! -name "*_impl*.go")
endef

define mockgen_service
	$(shell mockgen -source ./src/service/$1 -destination ./src/service/mock/$1)
endef

define get_repository_list
	$(shell	find ./src/domain/repository -maxdepth 1 -type f )
endef

define mockgen_repository
	$(shell mockgen -source ./src/domain/repository/$1 -destination ./src/domain/repository/mock/$1)
endef
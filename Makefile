
PROJECT_ID = abyssparanoia/rapid-go

PROJECT_DIR = /go/src/github.com/${PROJECT_ID}

#SERVICE_LIST = $(find ${PROJECT_DIR}/src/service/*.go -maxdepth 1 -type f ! -name "*_impl.go")

init:
	@echo Initialize rapid-go now......
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	@echo Initialize rapid-go completed!!!!

build:
	docker-compose build

start:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs api

mockgen_task:
	$(eval SERVICE_LIST := $(call get_service_list))	
	$(foreach file, $(SERVICE_LIST), $(call mockgen_service,$(shell basename $(file))))
	$(eval REPOSITORY_LIST := $(call get_repository_list))
	$(foreach file, $(REPOSITORY_LIST), $(call mockgen_repository,$(shell basename $(file))))

define get_service_list
	$(shell	docker-compose exec api find $(PROJECT_DIR)/src/service/ -maxdepth 1 -type f ! -name "*_impl*.go")
endef

define mockgen_service
	$(shell docker-compose exec api mockgen -source $(PROJECT_DIR)/src/service/$1 -destination $(PROJECT_DIR)/src/service/mock/$1)
endef

define get_repository_list
	$(shell	docker-compose exec api find $(PROJECT_DIR)/src/domain/repository -maxdepth 1 -type f )
endef

define mockgen_repository
	$(shell docker-compose exec api mockgen -source $(PROJECT_DIR)/src/domain/repository/$1 -destination $(PROJECT_DIR)/src/domain/repository/mock/$1)
endef
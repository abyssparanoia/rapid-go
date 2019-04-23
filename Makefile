
PROJECT_ID = abyssparanoia/rapid-go

PROJECT_DIR = /go/src/github.com/${PROJECT_ID}

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

mockgen:
	docker-compose exec api mockgen -source ${PROJECT_DIR}/src/${layer}/${interface}.go -destination ${PROJECT_DIR}/src/${layer}/mock/${interface}.go
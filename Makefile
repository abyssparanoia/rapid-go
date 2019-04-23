
init:
	@echo Initialize rapid-go now......
	go get -u github.com/golang/dep/cmd/dep
	dep ensure
	@echo Initialize rapid-go completed!!!!

build:
	docker-compose build

start:
	docker-compose up -d

stop:
	docker-compose down

logs:
	docker-compose logs api
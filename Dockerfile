FROM golang:1.11

COPY . /go/src/github.com/abyssparanoia/rapid-go/
WORKDIR /go/src/github.com/abyssparanoia/rapid-go/

RUN apt-get update -y \
    && apt-get install git -y \
    && go get -u github.com/codegangsta/gin \
    && go get -u github.com/golang/dep/cmd/dep \
    && go get -u github.com/golang/mock/gomock \
    && go install github.com/golang/mock/mockgen \
    && dep ensure

CMD gin -i run main.go routing.go dependency.go

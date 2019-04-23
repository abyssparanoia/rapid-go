for file in `\find ./src/service/*.go -maxdepth 1 -type f ! -name "*_impl.go" `; do
    docker-compose exec api mockgen -source ${PROJECT_DIR}/src/service/${file}.go -destination ${PROJECT_DIR}/src/service/mock/${file}.go
done
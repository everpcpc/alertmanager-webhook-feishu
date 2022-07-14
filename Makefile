fmt:
	go fmt ./...
run:fmt
	go run main.go server -c config.yml -e -v
build:
	goreleaser release --snapshot
docker_build:
	docker buildx build -t everpcpc/alertmanager-webhook-feishu --platform linux/amd64,linux/arm64 --push .

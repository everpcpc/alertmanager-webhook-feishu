fmt:
	go fmt ./...
run:fmt
	go run main.go server -c config.yml -e -v
build:
	goreleaser release --snapshot
docker_build:
	docker build -t everpcpc/alertmanager-webhook-feishu .
docker_push:docker_build
	docker push everpcpc/alertmanager-webhook-feishu

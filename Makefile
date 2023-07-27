start:
	go run ./cmd/

mock:
	mockgen -source=./service/manager.go -destination=./service/mock/mock.go *
	
test:
	go test -v -cover ./...

redis:
	docker run --rm --name redis -d -p 6379:6379 redis

.PHONY: start mock test
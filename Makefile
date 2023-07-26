start:
	go run ./cmd/

mock:
	mockgen -source=./service/manager.go -destination=./service/mock/mock.go *
	
test:
	go test -v -cover ./...

.PHONY: start mock test
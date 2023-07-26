start:
	go run ./cmd/

mock:
	mockgen -source=./service/manager.go -destination=./service/mock/mock.go *
	

.PHONY: start mock
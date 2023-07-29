start:
	rm ./log/*
	touch ./log/logs.txt
	go run ./cmd/

mock:
	mockgen -source=./service/manager.go -destination=./service/mock/mock.go *
	mockgen -source=./repository/manager.go -destination=./repository/mock/mock.go *
	
test:
	go test -v -cover ./...

redis:
	docker run --rm --name redis -d -p 6379:6379 redis

postgres:
	docker run --rm --name psg -d -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWROD=secret123 -e POSTGRES_DB=tsarka postgres

postgres-stop:
	docker exec psg pg_dump -U tsarka tsarka > tsarka_dump.sql
	docker stop psg

.PHONY: start mock test redis
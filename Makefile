.SILENT:

run:
	go run cmd/app/main.go 

migrate:
	migrate -path ./schema -database "postgres://postgres:qwerty@192.168.56.1:5436/postgres?sslmode=disable" up
migrate-down:
	migrate -path ./schema -database "postgres://postgres:qwerty@192.168.56.1:5436/postgres?sslmode=disable" down 1

migrate-drop:
	migrate -path ./schema -database "postgres://postgres:qwerty@192.168.56.1:5436/postgres?sslmode=disable" drop
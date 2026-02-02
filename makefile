init:
	touch ./db.sqlite
	go install gorm.io/cli/gorm@latest
	go run ./cmd/gorm-gen/main.go

run:
	go run .

build:
	go build -o web-demo

migrate:
	~/go/bin/gorm gen -i ./ -o ./generated

tidy:
	go mod tidy
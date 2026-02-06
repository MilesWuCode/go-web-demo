init:
	touch ./db.sqlite
	go install gorm.io/cli/gorm@latest
	go run ./cmd/gorm-gen/main.go

run:
	go run .

build:
	go build -o web-demo

tidy:
	go mod tidy
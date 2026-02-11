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

rustfs:
	docker run -d \
	--name rustfs \
	-p 9000:9000 \
	-p 9001:9001 \
	-v $(pwd)/rustfs/data:/data \
	-v $(pwd)/rustfs/logs:/logs \
	-e RUSTFS_ACCESS_KEY=rustfsadmin \
  	-e RUSTFS_SECRET_KEY=rustfsadmin \
	-t rustfs/rustfs:latest

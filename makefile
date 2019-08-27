docker-build:
	docker build --tag liron-navon/code-runner .

docker-run:
	docker run -p 8080:80 -it liron-navon/code-runner

test:
	go test -cover ./...

start:
	go run main.go

build:
	go build
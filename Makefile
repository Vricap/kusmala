run:
	go run main.go

build:
	go build -o kusmala .
	./kusmala

test:
	go test ./lexer

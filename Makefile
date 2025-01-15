build:
	go build -o ./bin/kusmala .
	# ./bin/kusmala

run:
	go run main.go

test:
	go test ./lexer
	go test ./parser
	go test ./evaluator

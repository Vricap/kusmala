build:
	go build -o ./bin/kusmala .
	# ./bin/kusmala

build_linux_64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/kusmala . # static linked binary for linux 64

build_win_64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/kusmala.exe . # static linked binary for windows 64

run:
	go run main.go

test:
	go test ./lexer
	go test ./parser
	go test ./evaluator

mod:
	go build -o ./bin/kusmala .
	./bin/kusmala test.km
	./bin/kusmala 1.km
	./bin/kusmala 2.km
	./bin/kusmala 3.km
	./bin/kusmala 4.km
	./bin/kusmala ./contoh/faktorial.km
	./bin/kusmala ./contoh/fibonacci.km
	./bin/kusmala ./contoh/kompleks_jika.km
	./bin/kusmala ./contoh/fungsi_dan_jika.km
	./bin/kusmala ./contoh/loop.km

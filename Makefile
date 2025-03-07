build:
	go build -o ./bin/kusmala main.go
	# ./bin/kusmala

build_linux_64:
	# static linked binary for linux 64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/kusmala_linux_amd64 main.go

build_win_64:
	# static linked binary for windows 64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/kusmala_windows_amd64.exe main.go

run:
	go run main.go

install:
	 # go for some reason doesn't allow to specify the output binary name in 'go install'
	go install .

test:
	go test ./lexer
	go test ./parser
	go test ./evaluator

mod:
	go build -o ./bin/kusmala main.go
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

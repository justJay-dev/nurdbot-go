dev:
	nodemon --ignore ./database/ --exec "go run" src/main.go

run:
	go run src/main.go

build:
	go build -o bin/main src/main.go

clean:
	rm -rf bin/*

install:
	go get
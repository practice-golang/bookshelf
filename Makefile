build:
	go build -o bin/

clean:
	rm -rf ./bin

test:
	go test bookshelf book db

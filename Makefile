build:
	go build -o bin/

clean:
	rm -f bin/bookshelf*
	rm -f bin/*log

test:
	go test bookshelf book db

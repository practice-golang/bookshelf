build:
	go build -o bin/

modvendor:
	go build -mod=vendor -o bin/

clean:
	rm -rf ./bin

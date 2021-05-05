build:
	go build -o bin/

vendor:
	go build -mod=vendor -o bin/

clean:
	rm -rf ./bin

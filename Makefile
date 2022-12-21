all:
	go build .

clean:
	go clean
	rm bin/*

bin:
	GOOS=windows GOARCH=amd64 go build -o bin/asciiImage-amd64.exe .
	GOOS=darwin GOARCH=amd64 go build -o bin/asciiImage-amd64-darwin . 
	GOOS=linux GOARCH=amd64 go build -o bin/asciiImage-amd64-linux .

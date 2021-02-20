install:
	go install

windows:
	env GOOS=windows GOARCH=amd64 go build

build:
	go build
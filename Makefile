run:
	go run main.go

build:
	GOOS=windows GOARCH=amd64 go build -o dist/JavaFinder-Windows-64bit/javafinder.exe main.go
	GOOS=linux GOARCH=amd64 go build -o dist/JavaFinder-Linux-64bit/javafinder main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/JavaFinder-MacOS-64bit/javafinder main.go
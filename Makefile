run:
	go run main.go

build:
	GOOS=windows GOARCH=amd64 go build -o dist/OJDMCollector-Windows-amd64bit/ojdm-collector.exe main.go
	GOOS=linux GOARCH=amd64 go build -o dist/OJDMCollector-Linux-amd64bit/ojdm-collector main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/OJDMCollector-MacOS-amd64bit/ojdm-collector main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/OJDMCollector-MacOS-arm64bit/ojdm-collector main.go
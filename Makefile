# use the vendor/ subdir which holds the vendored apache thrift go library, version
# the vendored thrift is commit fa0796d33208eadafb6f42964c8ef29d7751bfc2 on 1.0.0-dev,
# last commit there is Fri Oct 16 21:33:39 2015 +0200, from https://github.com/apache/thrift

test:
	go test --cover ./...

all:
	make windows
	make linux
	make mac

build:
	make all
	cp ./.env ./bin/.env
	echo "Build Success!"

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./bin/x_win_x86.exe main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/x_win_x64.exe main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/x_linux_x86 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/x_linux_x64 main.go

mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o ./bin/x_osx_x86 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/x_osx_64 main.go
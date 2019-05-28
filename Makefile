# use the vendor/ subdir which holds the vendored apache thrift go library, version
# the vendored thrift is commit fa0796d33208eadafb6f42964c8ef29d7751bfc2 on 1.0.0-dev,
# last commit there is Fri Oct 16 21:33:39 2015 +0200, from https://github.com/apache/thrift

test:
	go test --cover -covermode=count -coverprofile=coverage.out ./...

all:
	make user_win
	make user_linux
	make user_mac
	make admin_win
	make admin_linux
	make admin_mac
	make message_queue_win
	make message_queue_linux
	make message_queue_mac

build:
	make all
	cp ./.env ./bin/.env
	echo "Build Success!"

user_win:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./bin/user_win_x86.exe ./cmd/user/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/user_win_x64.exe ./cmd/user/main.go

user_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/user_linux_x86 ./cmd/user/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/user_linux_x64 ./cmd/user/main.go

user_mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o ./bin/user_osx_x86 ./cmd/user/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/user_osx_64 ./cmd/user/main.go
	
admin_win:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./bin/admin_win_x86.exe ./cmd/admin/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/admin_win_x64.exe ./cmd/admin/main.go

admin_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/admin_linux_x86 ./cmd/admin/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/admin_linux_x64 ./cmd/admin/main.go

admin_mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o ./bin/admin_osx_x86 ./cmd/admin/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/admin_osx_64 ./cmd/admin/main.go

message_queue_win:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ./bin/message_queue_win_x86.exe ./cmd/message_queue/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/message_queue_win_x64.exe ./cmd/message_queue/main.go

message_queue_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o ./bin/message_queue_linux_x86 ./cmd/message_queue/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/message_queue_linux_x64 ./cmd/message_queue/main.go

message_queue_mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o ./bin/message_queue_osx_x86 ./cmd/message_queue/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/message_queue_osx_64 ./cmd/message_queue/main.go
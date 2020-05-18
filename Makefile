# use the vendor/ subdir which holds the vendored apache thrift go library, version
# the vendored thrift is commit fa0796d33208eadafb6f42964c8ef29d7751bfc2 on 1.0.0-dev,
# last commit there is Fri Oct 16 21:33:39 2015 +0200, from https://github.com/apache/thrift

test:
	GO_TESTING=1 go test -mod=vendor --cover -covermode=count -coverprofile=coverage.out ./...
	bash ./scripts/clean.sh

build:
	bash ./scripts/build.sh
	echo "Build Success!"

clean:
	bash ./scripts/clean.sh
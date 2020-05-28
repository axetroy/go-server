# Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
test:
	GO_TESTING=1 go test -mod=vendor --cover -covermode=count -coverprofile=coverage.out ./...

build:
	bash ./scripts/build.sh admin
	bash ./scripts/build.sh user
	bash ./scripts/build.sh resource
	bash ./scripts/build.sh message_queue
	bash ./scripts/build.sh scheduled
	echo "Build Success!"

clean:
	bash ./scripts/clean.sh

# deploy app via s4 see detail in https://github.com/axetroy/s4
deploy:
	s4
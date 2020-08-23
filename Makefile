# Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
test:
	GO_TESTING=1 go test -timeout=30m -mod=vendor --cover -covermode=count -coverprofile=coverage.out ./...

build:
	bash ./scripts/build.sh admin
	bash ./scripts/build.sh user
	bash ./scripts/build.sh resource
	bash ./scripts/build.sh message_queue
	bash ./scripts/build.sh scheduled
	bash ./scripts/build.sh customer_service
	echo "Build Success!"

clean:
	bash ./scripts/clean.sh

# deploy app via s4 see detail in https://github.com/axetroy/s4
deploy:
	s4

generate-static:
	pkger -include /internal/service/area/external -o ./internal/service/area/external_pkged
	pkger -include /internal/app/customer_service/views -o ./internal/app/customer_service/views_pkged

lint:
	golangci-lint run ./... -v

format:
	go fmt ./...
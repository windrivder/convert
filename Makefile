.PHONY: init
# init
init:
	apt-get install musl-tools -y

.PHONY: build
# build
build:
	CC=musl-gcc CXX=musl-g++ CGO_ENABLED=1 go build -tags "musl" -o convert .

.PHONY: buildx
# buildx
buildx:
	docker build --network=host -t convert:v1.0 .

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: install-deps
install-deps:
	brew install go
	brew install protobuf
	brew install libpq
	brew link --force libpq
	brew install bufbuild/buf/buf


.PHONY: fmt-proto
fmt-proto:
	buf format proto --diff --exit-code

.PHONY: gen-code
gen-code:
	buf generate --timeout 5m --template buf.gen.go.yaml

.PHONY: start-services
start-services:
	docker compose -f .docker/docker-compose.yaml up

tidy:
	go mod tidy

codegen-api:
	oapi-codegen \
	-generate gin,strict-server,types,spec \
	-package http -o internal/api/http/http.gen.go docs/oapi/api.yaml

codegen-mocks:
	mockery

codegen-di:
	wire github.com/SergeiGD/golang-template/internal/di

tests:
	go test -count=1 ./...

run:
	docker compose -f docker-compose.yaml up --build
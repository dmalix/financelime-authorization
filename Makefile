dep:
	go mod tidy
	go mod vendor

check:
	go fmt ./...
	go vet ./...
	go test ./...

build:
	scripts/build.sh

swagger:
	swag init
	goswagger generate markdown --spec=docs/swagger.json --output=docs/api.md
	mv docs/api.md docs/API.md

staging:
	scripts/staging.sh

all: dep check build swagger


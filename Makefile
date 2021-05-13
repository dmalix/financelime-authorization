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
	swag init --output contract
	goswagger generate markdown --spec=contract/swagger.json --output=contract/api.md
	mv contract/api.md contract/API.md

staging:
	scripts/staging.sh

all: dep check build swagger


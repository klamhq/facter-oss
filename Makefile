OUTPUT:= bin
BIN := facter-oss
BIN_MAC := facter-oss
BUILD_ENV := GOOS=linux GOARCH=amd64
BUILD_FLAGS := -tags netgo -ldflags '-w'
VERBOSE := -x -v

.PHONY: all

mkdocs-install:
	python3 -m venv venv/
	source venv/bin/activate
	pip3 install mkdocs
	pip3 install mkdocs-material
	python3 -m pip freeze > requirements.txt
	

build:
	@mkdir -p ${OUTPUT}
	${BUILD_ENV} go build -mod vendor ${BUILD_FLAGS} -o ${OUTPUT}/ ./...
	@chmod +x ${OUTPUT}/${BIN}

buildMac:
	@mkdir -p ${OUTPUT}
	GOARCH=amd64 go build -mod vendor ${BUILD_FLAGS} -o ${OUTPUT}/ ./...
	@chmod +x ${OUTPUT}/${BIN_MAC}

test:
	go test -v -cover -race ./... -coverprofile=./tests/coverage.out && go tool cover -html=./tests/coverage.out -o ./tests/cover.html

sudoTest:
	sudo go test -v -cover -race ./... -coverprofile=./tests/coverage.out && go tool cover -html=./tests/coverage.out -o ./tests/cover.html

clean:
	@rm -fr ${CURDIR}/bin ${CURDIR}/vendor ${CURDIR}/data

install_dependencies:
	go mod vendor

check_compliance: build
	go run tools/compliant/main.go

compress: build
	upx -8 ${OUTPUT}/${BIN}

release: build compress

dockerRocky:
	docker build -t rockylinux:test -f docker/Dockerfile .

dockerRocky8:
	docker build -t rockylinux8:test  -f  docker/Dockerfile-rocky8 .

dockerUbuntu:
	docker build -t ubuntu:test -f  docker/Dockerfile-ubuntu .

test-docker-rocky8:
	docker run --name rockylinux8 -h rockylinux8 --rm -v $(CURDIR)/:/tmp/facter rockylinux8:test

test-docker-ubuntu:
	docker run --name ubuntu -h ubuntu --rm -v $(CURDIR)/:/tmp/facter ubuntu:test

test-docker-rocky:
	docker run --name rockylinux -h rockylinux --rm -v $(CURDIR)/:/tmp/facter rockylinux:test

all: test build

golint:
	golangci-lint run

profile-clean:
	-sudo rm -f cpu-perf mem-perf

profile: profile-clean
	go build -o ${OUTPUT}/${BIN}-perf main.go
	sudo ${OUTPUT}/${BIN}-perf --config ./configs/config.yml
	go tool pprof -http 127.0.0.1:9002 cpu-perf

decode-protobuf:
	 protoc --decode klamhq.rpc.facter.v1.InventoryRequest -I ./vendor/github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1 ./vendor/github.com/klamhq/facter-schema/proto/klamhq/rpc/facter/v1/service.proto < /tmp/export.iya
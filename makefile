# This is how we want to name the binary output
BINARY=lighthouse

# These are the values we want to pass for Version and BuildTime
VERSION=0.0.1-SNAPSHOT
BUILD_TIME=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

all: makebin test
	go build ${LDFLAGS} -o ${BINARY} 

test: makebin
	go test -v ./...

makebin:
	go get -u github.com/jteeuwen/go-bindata/...
	date "+%Y%m%d%H%M%S" > static/CREATE_TIME
	go-bindata -o bindata/bindata.go -pkg bindata -ignore /\\..* -prefix static/ static/...

run: all
	open "https://127.0.0.1:8443" && ./${BINARY}

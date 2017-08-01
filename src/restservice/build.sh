export GOPATH=$(pwd)
export GOROOT=/usr/bin/

echo "Fetching imports..."
go get "github.com/julienschmidt/httprouter"

echo "Building golang code..."
go build main.go

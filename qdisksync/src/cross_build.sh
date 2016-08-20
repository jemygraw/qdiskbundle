DIR=$(cd ../; pwd)
export GOPATH=$GOPATH:$DIR
GOOS=linux GOARCH=amd64 go build -o qdisksync main.go
GOOS=windows GOARCH=amd64 go build -o qdisksync.exe main.go
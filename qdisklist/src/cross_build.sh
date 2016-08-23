DIR=$(cd ../; pwd)
export GOPATH=$GOPATH:$DIR
GOOS=linux GOARCH=amd64 go build -o qdisklist main.go
GOOS=windows GOARCH=386 go build -o qdisklist.exe main.go

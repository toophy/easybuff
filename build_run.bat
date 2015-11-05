rm -rf bin/proto_b
go build -o ./bin/easybuff.exe main.go
cd ./bin/
easybuff.exe $*
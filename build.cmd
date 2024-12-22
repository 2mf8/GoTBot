SET GOOS=linux
SET GOARCH=amd64
go build -o QQHelperAMD64 .

SET GOOS=linux
SET GOARCH=arm64
go build -o QQHelperARM64 .
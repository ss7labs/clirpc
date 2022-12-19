go build -ldflags="-s -w" -v cmd/rpctr.go

if [ "$?" -ne 0 ]
then
	exit
fi

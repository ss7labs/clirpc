go build -ldflags="-s -w" -v cmd/rpctr.go

if [ "$?" -ne 0 ]
then
	exit
fi
ssh -i ~/key-store/id_rsa_bras jetb@10.19.176.55 'pkill -f "rpctr"'
scp -i ~/key-store/id_rsa_bras rpctr jetb@10.19.176.55:
ssh -i ~/key-store/id_rsa_bras jetb@10.19.176.55 '/home/jetb/run.sh < /dev/null > /tmp/rpctr.log 2>&1 &'

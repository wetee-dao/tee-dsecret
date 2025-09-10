# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/node1

export SIDE_CHAIN_PORT=41000
export GQL_PORT=41005
export PASSWORD=123456
export CHAIN_ADDR=ws://192.168.110.205:30002/ws

ego-go build -o dsecret ../../main.go
ego sign dsecret

rm nohup.out

ego sign dsecret && ego run dsecret
# nohup ego run dsecret &
# ./dsecret
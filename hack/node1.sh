# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/node1

export PEER_PK=080112406bce93c01f4b51287b01e55565cf7933cb624b25d478e003ca23446bc3ef83b9d0380163fd5c55a0474b95709da5b31d386da0313bb69bd635618f5cb80f1dde
export TCP_PORT=31000
export UDP_PORT=31000
export GQL_PORT=31001
export PASSWORD=123456
export CHAIN_ADDR=ws://xiaobai.asyou.me:30002/ws

ego-go build -o dsecret ../../main.go
ego sign dsecret
ego run dsecret
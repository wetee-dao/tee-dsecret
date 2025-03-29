# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/node2

export PEER_PK=08011240c137a203132c2fb66f13de24f4e1db4177daa5d334c51afeb3aa195db414fea8a2babbb311378d1a707a940a171947d80202fdc1799923e9b045393f58d18472
export TCP_PORT=31002
export UDP_PORT=31002
export GQL_PORT=31003
export PASSWORD=123456
export CHAIN_ADDR=ws://paseo.asyou.me/ws

ego-go build -o dsecret ../../main.go
ego sign dsecret

rm nohup.out

nohup ego run dsecret &
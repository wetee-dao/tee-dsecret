# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/node2

export PEER_PK='casino fitness lens stable viable alter gossip year game suspect zero surface'
export TCP_PORT=32001
export UDP_PORT=32001
export BOOT_PEERS='/ip4/0.0.0.0/tcp/32000/p2p/12D3KooWMiR8pUZBzf7aY9cnTrUWbTLYJTJCv82mMYZtf1XZTuYt'

go run ../../main.go

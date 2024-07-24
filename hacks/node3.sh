# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/node3

export PEER_PK='trumpet news diet produce faith measure rhythm cry pink noodle saddle glad'
export TCP_PORT=32002
export UDP_PORT=32002
export BOOT_PEERS='/ip4/0.0.0.0/tcp/32000/p2p/12D3KooWMiR8pUZBzf7aY9cnTrUWbTLYJTJCv82mMYZtf1XZTuYt'

go run ../../main.go

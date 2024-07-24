# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/node1

export PEER_PK='broken lounge citizen various summer opera sleep there brother rely voyage cash'
export TCP_PORT=32000
export UDP_PORT=32000

go run ../../main.go

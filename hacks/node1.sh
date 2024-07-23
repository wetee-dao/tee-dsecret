# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/../

export PEER_PK='broken lounge citizen various summer opera sleep there brother rely voyage cash'
export TCP_PORT=31000
export UDP_PORT=31000
export NODES='[{"id":"d0380163fd5c55a0474b95709da5b31d386da0313bb69bd635618f5cb80f1dde","address":"/ip4/127.0.0.1/tcp/31000"},{"id":"a2babbb311378d1a707a940a171947d80202fdc1799923e9b045393f58d18472","address":"/ip4/127.0.0.1/tcp/31001"},{"id":"99c19ae8c6ac65bb9bc48dcdd9d82e00ef7c8bfa341077ea261a8e0d2d7d7531","address":"/ip4/127.0.0.1/tcp/31002"}]'

go run main.go

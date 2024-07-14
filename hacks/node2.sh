# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/../

export PEER_PK=08011240c137a203132c2fb66f13de24f4e1db4177daa5d334c51afeb3aa195db414fea8a2babbb311378d1a707a940a171947d80202fdc1799923e9b045393f58d18472
export TCP_PORT=31001
export UDP_PORT=31001
export BOOT_PEERS=/ip4/127.0.0.1/tcp/31000/p2p/12D3KooWPqAW35BWBWk9N6MwYwoCdzk4TKVKidhhoNpxwtekPsNM
export NODES='[{"id":"d0380163fd5c55a0474b95709da5b31d386da0313bb69bd635618f5cb80f1dde","address":"/ip4/127.0.0.1/tcp/31000"},{"id":"a2babbb311378d1a707a940a171947d80202fdc1799923e9b045393f58d18472","address":"/ip4/127.0.0.1/tcp/31001"},{"id":"99c19ae8c6ac65bb9bc48dcdd9d82e00ef7c8bfa341077ea261a8e0d2d7d7531","address":"/ip4/127.0.0.1/tcp/31002"}]'

go run main.go

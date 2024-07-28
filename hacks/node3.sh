# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/node3

export PEER_PK=08011240a2ccd0c8e266d32fbd65c1e790117bd55f8fcfdc7d203a57343275e3df0a98ce99c19ae8c6ac65bb9bc48dcdd9d82e00ef7c8bfa341077ea261a8e0d2d7d7531
export TCP_PORT=31002
export UDP_PORT=31002

go build -o dsecret ../../main.go
./dsecret

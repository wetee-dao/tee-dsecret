# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

ego-go build -o ./bin/dsecret ./main.go

export SIDE_CHAIN_PORT=30010
export GQL_PORT=30015
export CHAIN_ADDR=ws://192.168.110.205:30002/ws

cd ./bin/
ego sign dsecret && ego run dsecret
# ./dsecret
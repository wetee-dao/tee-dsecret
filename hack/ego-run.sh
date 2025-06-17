# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

export CHAIN_URI=wss://xiaobai.asyou.me:30001
ego-go build -o ./bin/dsecret ./main.go

export SIDE_CHAIN_PORT=31000
export GQL_PORT=31005
export PASSWORD=123456
export CHAIN_ADDR=ws://127.0.0.1:9944

cd ./bin/
# ego sign dsecret && ego run dsecret
./dsecret
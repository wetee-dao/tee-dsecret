# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

cd bin
rm -rf ./chain_data/BFT
rm -rf ./chain_data/wetee

cd $DIR/../hack/node1
rm -rf ./chain_data/BFT
rm -rf ./chain_data/wetee

cd $DIR/../hack/node2
rm -rf ./chain_data/BFT
rm -rf ./chain_data/wetee
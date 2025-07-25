# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd "$DIR/../../contract/"

cargo contract build --release --manifest-path contracts/Subnet/Cargo.toml

rm $DIR/contract_cache/subnet.*


cp target/ink/subnet/subnet.contract $DIR/contract_cache
cp target/ink/subnet/subnet.polkavm $DIR/contract_cache
cp target/ink/subnet/subnet.json $DIR/contract_cache

cd $DIR

cp ./contract_cache/subnet.json ../pkg/chains/contracts/


cd $DIR/../pkg/chains/contracts/
go-ink-gen -json subnet.json

go test -run ^TestSubnetUpdate$
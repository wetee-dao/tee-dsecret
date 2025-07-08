# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd "$DIR/../../contract/"

cargo contract build --release --manifest-path contracts/Pod/Cargo.toml
cargo contract build --release --manifest-path contracts/Subnet/Cargo.toml
cargo contract build --release --manifest-path contracts/Cloud/Cargo.toml

cp target/ink/cloud/cloud.contract $DIR/contract_cache
cp target/ink/cloud/cloud.polkavm $DIR/contract_cache
cp target/ink/cloud/cloud.json $DIR/contract_cache

cp target/ink/subnet/subnet.contract $DIR/contract_cache
cp target/ink/subnet/subnet.polkavm $DIR/contract_cache
cp target/ink/subnet/subnet.json $DIR/contract_cache

cp target/ink/pod/pod.contract $DIR/contract_cache
cp target/ink/pod/pod.polkavm $DIR/contract_cache
cp target/ink/pod/pod.json $DIR/contract_cache

cd $DIR

cp ./contract_cache/cloud.json ../pkg/chains/contracts/
cp ./contract_cache/subnet.json ../pkg/chains/contracts/
cp ./contract_cache/pod.json ../pkg/chains/contracts/


cd $DIR/../pkg/chains/contracts/
go-ink-gen -json cloud.json
go-ink-gen -json subnet.json
go-ink-gen -json pod.json

cd $DIR
# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd "$DIR/../../contract/"

cargo contract build --release --manifest-path contracts/Cloud/Cargo.toml

rm $DIR/contract_cache/cloud.*

cp target/ink/cloud/cloud.contract $DIR/contract_cache
cp target/ink/cloud/cloud.polkavm $DIR/contract_cache
cp target/ink/cloud/cloud.json $DIR/contract_cache

cd $DIR

cp ./contract_cache/cloud.json ../pkg/chains/contracts/


cd $DIR/../pkg/chains/contracts/
go-ink-gen -json cloud.json

go test -run ^TestCloudUpdate$
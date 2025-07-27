# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../pkg/model


protoc --proto_path=. --gogofast_out=. tx.proto

rm ../../../libos-entry/model/tx.pb.go
rm ../../../libos-entry/model/sgx_issue.go
cp tx.pb.go ../../../libos-entry/model/
cp sgx_issue.go ../../../libos-entry/model/
cp -r protoio ../../../libos-entry/model/
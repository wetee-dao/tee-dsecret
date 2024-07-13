# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/../

export PEER_PK=0801124085e1ef6c9a42d90ca0fa02feac23bd59482122ac31d06586406a28466ea169fdd9dc01ca76a12de425cbb6cef37f286e5937759a5c2b59f5d53a6c3449fea113
export TCP_PORT=31002
export UDP_PORT=31002
export BOOT_PEERS=/ip4/127.0.0.1/tcp/31000/p2p/12D3KooWEA5ycwyyRKk3vgnRKErtqCVBqvk4pdGUSNDTesYDA95E
export PKG_PK=2ae230cd9de7e2e7ee51559aa90e40395ddc699cc1268c6ab1feb184fbc38508
export PKG_PUBS=f84227f038b5a4d6b7d40291c79fce151cf635e9187fee1a1b64021b8386017c_ba6eaf20408989f46a6ce7b255f9173a16031cc1ebc44ae227550158c7c6564a_d94845a8ab8d7e07354182e31a7fb2478d95dd6700223afb62414185de657b49

go run main.go

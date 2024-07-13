# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/../

export PEER_PK=0801124004c88e4c7fc0dddfe064f945e21f826214bc408ab56ce6d4e6a7782e2f45a7443a63561b4b9db9708644eda677fa1532b5f989bac9456de47cc499a7a4947dc5
export TCP_PORT=31001
export UDP_PORT=31001
export BOOT_PEERS=/ip4/127.0.0.1/tcp/31000/p2p/12D3KooWEA5ycwyyRKk3vgnRKErtqCVBqvk4pdGUSNDTesYDA95E
export PKG_PK=3c11a61ef0f73ede226f4719d811697968b5294d00ae43d852e70a0c610cee01
export PKG_PUBS=f84227f038b5a4d6b7d40291c79fce151cf635e9187fee1a1b64021b8386017c_ba6eaf20408989f46a6ce7b255f9173a16031cc1ebc44ae227550158c7c6564a_d94845a8ab8d7e07354182e31a7fb2478d95dd6700223afb62414185de657b49

go run main.go

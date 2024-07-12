# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

cd $DIR/../

export PEER_PK=08011240de21840c7b9084b909730d73ed17393309b727cd35887a707d149828c5efa788407c0e47055c0695aeaaaaeeed8dcb55484fe8178ffd963ea98dcb7206c2eb9d
export TCP_PORT=61000
export UDP_PORT=61000
export SENDER=T
export PKG_PK=15bb2b22e24979a89555c67648d2b1994a0309a79defd650f256ea3f36b54502
export PKG_PUBS=f84227f038b5a4d6b7d40291c79fce151cf635e9187fee1a1b64021b8386017c_ba6eaf20408989f46a6ce7b255f9173a16031cc1ebc44ae227550158c7c6564a_d94845a8ab8d7e07354182e31a7fb2478d95dd6700223afb62414185de657b49

go run main.go

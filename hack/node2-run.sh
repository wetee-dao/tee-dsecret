# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"

# Get image version from .version file
img=$(cat $DIR/.version)

cd $DIR/node2

# Ensure chain_data directory exists
mkdir -p $DIR/node2/chain_data

# Run dsecret using Docker
docker run --name dsecret-2 \
  --rm \
  -p 30030:30130 \
  -p 30035:30135 \
  -e SIDE_CHAIN_PORT=30130 \
  -e GQL_PORT=30135 \
  -e CHAIN_ADDR=ws://192.168.110.205:30002/ws \
  -v $DIR/node2/chain_data:/chain_data \
  --device /dev/sgx_enclave \
  --device /dev/sgx_provision \
  $img
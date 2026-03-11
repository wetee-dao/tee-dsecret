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

cd $DIR/node0

# Ensure chain_data directory exists
mkdir -p $DIR/node0/chain_data

# Run dsecret using Docker
docker run --name dsecret-0 \
  --rm \
  -p 30010:30110 \
  -p 30011:30111 \
  -p 30015:30115 \
  -e SIDE_CHAIN_PORT=30110 \
  -e SIDE_CHAIN_RPC_PORT=30111 \
  -e GQL_PORT=30115 \
  -e CHAIN_ADDR=ws://192.168.110.205:30002/ws \
  -v $DIR/node0/chain_data:/chain_data \
  --device /dev/sgx_enclave \
  --device /dev/sgx_provision \
  $img
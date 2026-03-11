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

cd $DIR/node1

# Ensure chain_data directory exists
mkdir -p $DIR/node1/chain_data

# Run dsecret using Docker
docker run --name dsecret-1 \
  --rm \
  -p 30020:30120 \
  -p 30025:30125 \
  -e SIDE_CHAIN_PORT=30120 \
  -e GQL_PORT=30125 \
  -e CHAIN_ADDR=ws://192.168.110.205:30002/ws \
  -v $DIR/node1/chain_data:/chain_data \
  --device /dev/sgx_enclave \
  --device /dev/sgx_provision \
  $img

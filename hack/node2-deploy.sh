# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

img=$(cat ./hack/.version)

export DSECRET_IMAGE=$img
export DSECRET_DIR=/home/wetee/work/wetee/tee-dsecret/hack/node2/chain_data

export SIDE_CHAIN_PORT=30030
export GQL_PORT=30035
export CHAIN_ADDR=ws://192.168.110.205:30002/ws


echo '' > ./hack/k8s.yaml
envsubst < ./hack/k8s-temp.yaml > ./hack/k8s.yaml

# 部署镜像
kubectl delete -f ./hack/k8s.yaml
kubectl create -f ./hack/k8s.yaml
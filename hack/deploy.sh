# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

tag=`date "+%Y-%m-%d-%H_%M"`

# 构建镜像
ego-go build -o ./bin/dsecret ./server.go

docker build -t registry.cn-hangzhou.aliyuncs.com/wetee_dao/dsecret:$tag .

docker push registry.cn-hangzhou.aliyuncs.com/wetee_dao/dsecret:$tag

export WETEE_INDEXER_IMAGE=registry.cn-hangzhou.aliyuncs.com/wetee_dao/dsecret:$tag
echo '' > ./hack/dsecret.yaml
envsubst < ./hack/dsecret-temp.yaml > ./hack/dsecret.yaml

# # 部署镜像
kubectl create -f ./hack/dsecret.yaml

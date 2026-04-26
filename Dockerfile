FROM wetee/ego-ubuntu-24-04:1.8.1
WORKDIR /

RUN mkdir -p /chain_data

ADD hack/build/dsecret  /

CMD ["/bin/sh", "-c" ,"ego run dsecret"]
FROM wetee/ego-ubuntu-24-04:1.7.2
WORKDIR /

ADD hack/build/dsecret  /

RUN mkdir -p /chain_data

CMD ["/bin/sh", "-c" ,"ego run dsecret"]
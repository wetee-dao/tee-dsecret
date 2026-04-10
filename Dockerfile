FROM wetee/ego-ubuntu-24-04:1.7.2
WORKDIR /

RUN mkdir -p /chain_data

ADD hack/build/dsecret  /

CMD ["/bin/sh", "-c" ,"ego run dsecret"]
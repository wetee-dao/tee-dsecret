FROM wetee/ego-ubuntu-deploy:22.04:22.04
WORKDIR /

ADD bin/manager  /

RUN mkdir -p /opt/wetee-worker

EXPOSE 8880 8883 

CMD ["/bin/sh", "-c" ,"ego run manager"]
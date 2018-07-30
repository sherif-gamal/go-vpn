FROM golang

COPY locale.gen /etc/
COPY .tmux.conf ~
RUN go get github.com/songgao/packets/ethernet \
    && go get github.com/songgao/water \
    && go get github.com/songgao/water/waterutil \
    && go get golang.org/x/net/ipv4 \
    && apt-get update \
    && apt-get install -y locales vim tmux net-tools
RUN export PATH="/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

WORKDIR /go/src/github.com/sherif-gamal/go-vpn
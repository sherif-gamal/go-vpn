This is a minimal vpn written in golang that creates a tap interface to intercepts packets sent
between hosts, encrypts them then sends them to the other host.

This code has been tested on "Debian GNU/Linux 9 (stretch)".

To test on the same environment build the docker image

> docker build ./

then run the docker image with this command

> docker run -it --cap-add=NET_ADMIN --device=/dev/net/tun -v \`pwd\`:/go/src/github.com/sherif-gamal/go-vpn {imageid}

build the code

> go build

then you can run the code as follows

> ./go-vpn -local 10.0.1.1 -remote 10.0.2.1 -lport 9999 -rport 8888 -iface i1

and in another session (on the same container):

> ./go-vpn -local 10.0.2.1 -remote 10.0.1.1 -lport 8888 -rport 9999 iface i2

now if you ping any ip in a subnet of either interface you should see the log messages of both instances

> ping 10.0.1.2

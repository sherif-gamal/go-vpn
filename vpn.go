package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"

	"github.com/songgao/water"
)

const BUFFERSIZE = 2000

var c Crypto

var (
	localIP  = flag.String("local", "10.0.0.1", "Local tun interface IP/MASK like 192.168.3.3/24")
	lport    = flag.String("lport", "9999", "local port")
	remoteIP = flag.String("remote", "10.0.1.2", "Remote server (external) IP like 8.8.8.8")
	rport    = flag.String("rport", "9998", "Remote port")
	iname    = flag.String("iface", "O_O", "Interafce name")
)

func runCmd(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		log.Fatalln("Error running ", command, ": ", err)
	}
}

func initInterface() *water.Interface {
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = *iname

	iface, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	runCmd("ip", "addr", "add", fmt.Sprintf("%v/24", *localIP), "dev", *iname)
	runCmd("ip", "link", "set", "dev", *iname, "up")
	return iface
}

func listenUDP(conn *net.UDPConn, iface *water.Interface) {
	packet := make([]byte, BUFFERSIZE)
	for {
		log.Println("will wait to read from udp")
		n, remote, err := conn.ReadFromUDP(packet[:])
		log.Printf("Packet received from remote %+v", remote)
		if string(packet[:3]) != "ack" {
			conn.WriteToUDP([]byte("ack"), remote)
			if err != nil || n == 0 {
				log.Println("Error: ", err)
				continue
			}
		} else {
			log.Println("not sending ack")
		}
	}
	log.Println("here")
}

func listenToVirtualInterface(conn *net.UDPConn, iface *water.Interface) {
	packet := make([]byte, BUFFERSIZE)
	remoteAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%v", *remoteIP, *rport))

	for {
		log.Println("will wait to read from tap")
		n, err := iface.Read(packet[:])
		if err != nil {
			log.Fatal(err)
		}

		log.Println("read something")
		encrypted := c.Encrypt(packet[:])

		n, err = conn.WriteToUDP(encrypted[:], remoteAddr)
		if err != nil {
			log.Fatalf("error %+v", err)
		}
		log.Printf("sent %d bytes to %+v", n, remoteAddr)
	}
	log.Println("heresssss")
}

func server() *net.UDPConn {
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%v", *localIP, *lport))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("listening on %+v", conn.LocalAddr())
	return conn
}

func main() {
	flag.Parse()
	iface := initInterface()
	c.Init()
	conn := server()
	defer conn.Close()
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		listenUDP(conn, iface)
		wg.Done()
	}()
	go func() {
		listenToVirtualInterface(conn, iface)
		wg.Done()
	}()
	wg.Wait()
}

package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Proxy struct {
	cmap     map[string]*Connection // for each destination IP, there should be a connection
	readConn *net.UDPConn
	port     int
	proxyd   *ProxyData
	pmut     sync.Mutex
}

// proxy
// read data from source
// - get the source IP
//  -- from ProxyData get destination
//   -- from Prxy, get the connection/channel
//		?? if no connection, create new
//  -- send the data to channel
func StartProxy(pd *ProxyData, port int) (*Proxy, error) {

	// Set up Reciever
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	pudp, err := net.ListenUDP("udp", saddr)
	if err != nil {
		return nil, err
	}
	pxy := &Proxy{readConn: pudp, proxyd: pd, port: port, cmap: map[string]*Connection{}}
	log.Printf("Proxy reading on port %d\n", port)

	go reader(pxy)

	return pxy, nil
}

// Go routine which manages reads and sending to channels
func reader(pxy *Proxy) {
	var buffer [1500]byte
	for {
		// Read from server
		n, cliaddr, err := pxy.readConn.ReadFromUDP(buffer[0:])
		if err != nil {
			log.Printf("reader err: %+v", err)
			continue
		}
		// find destination for this client
		log.Print("Destination:" + cliaddr.IP.String())

		destip := pxy.proxyd.GetDest(cliaddr.IP.String())
		pxy.pmut.Lock()
		conp := pxy.cmap[destip]
		pxy.pmut.Unlock()

		if conp == nil {
			log.Printf("New Connection: %s", destip)
			daddr, err := net.ResolveUDPAddr("udp", destip)
			if err != nil {
				log.Printf("resolve [%s] err: %+v", destip, err)
				continue
			}
			dc := NewiDestConnection(daddr)
			if dc == nil {
				log.Printf("dest conn [%s] failed", destip)
				continue
			}
			pxy.cmap[destip] = dc // sender started
			conp = dc
		}

		// time to send
		dtos := buffer[0:n]
		conp.Worker <- dtos
	}
}

//func (p *ProxyData) AddSender(ip string, chan []byte) {

//
// -- Connection to Destination
//     and sender part
//

// Information maintained for each server connection
type Connection struct {
	Worker     chan []byte
	ServerConn *net.UDPConn // UDP connection to server
}

// Generate a new connection by opening a UDP connection to the server
func NewiDestConnection(srvAddr *net.UDPAddr) *Connection {
	conn := new(Connection)
	srvudp, err := net.DialUDP("udp", nil, srvAddr)
	if err != nil {
		log.Printf("!!! Server: %+v", srvAddr)
		log.Printf("!!! DialUDP failed: %+v", err)
		return nil
	}
	conn.ServerConn = srvudp
	conn.Worker = make(chan []byte, 1)

	// start worker in goroutine
	go sender(conn)

	return conn
}

// ip is the destination ip
func sender(c *Connection) {
	for {
		data := <-c.Worker

		// Relay it to client
		_, err := c.ServerConn.Write(data)
		if err != nil {
			log.Printf("Error sending to dest [%+v]: %+v", c.ServerConn, err)
		}
	}
}

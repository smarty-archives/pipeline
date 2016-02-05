package udp

import (
	"log"
	"net"
)

func BindSocket(bindAddress string) (*net.UDPConn, error) {
	if address, err := net.ResolveUDPAddr("udp", bindAddress); err != nil {
		log.Println("[WARN] UDP socket bind failure:", err)
		return nil, err
	} else if socket, err := net.ListenUDP("udp", address); err != nil {
		log.Println("[WARN] UDP socket bind failure:", err)
		return nil, err
	} else {
		log.Printf("[INFO] Listening for UDP datagrams on %s.\n", address)
		socket.SetReadBuffer(readBufferSize)
		return socket, nil
	}
}

package nozzle

import (
	"log"
	"net"
	"time"
)

type TCPSender struct {
	target string
	socket net.Conn
}

func NewTCPSender(target string) *TCPSender {
	return &TCPSender{target: target}
}

func (t *TCPSender) Send(message string) {
	for {
		if t.socket == nil {
			conn, err := net.Dial("tcp", t.target)
			if err != nil {
				retryDelay(err)
				continue
			}
			t.socket = conn
		}
		_, err := t.socket.Write([]byte(message))
		if err != nil {
			_ = t.socket.Close()
			t.socket = nil
			retryDelay(err)
			continue
		}
		break
	}
}

func retryDelay(err error) {
	log.Printf("Error sending to target: %s. Retrying", err)
	time.Sleep(1 * time.Second)
}

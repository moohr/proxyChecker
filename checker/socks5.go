package checker

import (
	"fmt"
	"io"
	"net"
	"time"
)

func IsSocks5(addr net.TCPAddr, semaphore chan struct{}, dialTimeout, readTimeout time.Duration, output io.Writer) {
	semaphore <- struct{}{}
	dialer := net.Dialer{Timeout: dialTimeout} // Set a timeout of 5 seconds
	conn, err := dialer.Dial("tcp", addr.String())
	if err != nil {
		<-semaphore
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(readTimeout))
	greeting := []byte{0x05, 0x01, 0x00}
	if _, err := conn.Write(greeting); err != nil {
		<-semaphore
		return
	}
	response := make([]byte, 2)
	if _, err := conn.Read(response); err != nil {
		<-semaphore
		return
	}
	if response[0] == 0x05 && response[1] != 0xFF {
		fmt.Fprintf(output, "socks5://%s\n", addr.String())
	}
	<-semaphore
	return
}

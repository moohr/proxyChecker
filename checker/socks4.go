package checker

import (
	"fmt"
	"io"
	"net"
	"time"
)

func IsSocks4(addr net.TCPAddr, semaphore chan struct{}, dialTimeout, readTimeout time.Duration, output io.Writer) {
	semaphore <- struct{}{}
	dialer := net.Dialer{Timeout: dialTimeout} // Set a timeout of 5 seconds
	conn, err := dialer.Dial("tcp", addr.String())
	if err != nil {
		<-semaphore
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(readTimeout))
	greeting := []byte{
		0x04,
		0x01,
		0x00, 0x50,
		0x01, 0x01, 0x01, 0x01,
		0x00}
	if _, err := conn.Write(greeting); err != nil {
		<-semaphore
		return
	}
	response := make([]byte, 8)
	if _, err := conn.Read(response); err != nil {
		<-semaphore
		return
	}
	if response[0] == 0x00 && response[1] == 0x5A {
		fmt.Fprintf(output, "socks4://%s\n", addr.String())
	}
	<-semaphore
	return
}

package checker

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func IsSocks4(addr net.TCPAddr, wg *sync.WaitGroup, semaphore chan struct{}) {
	semaphore <- struct{}{}
	defer wg.Done()
	dialer := net.Dialer{Timeout: 5 * time.Second} // Set a timeout of 5 seconds
	conn, err := dialer.Dial("tcp", addr.String())
	if err != nil {
		<-semaphore
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))
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
		fmt.Println("socks4://" + addr.String())
	}
	<-semaphore
	return
}

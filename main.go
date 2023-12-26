package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"proxyChecker/checker"
	"strconv"
	"strings"
	"sync"
)

const max = 5000

func getList(filename string) ([]net.TCPAddr, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var tcpAddrs []net.TCPAddr
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Split(strings.TrimSpace(line), ":")
		if len(temp) != 2 {
			continue
		}
		ip := temp[0]
		port, err := strconv.Atoi(temp[1])
		if err != nil {
			continue
		}
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
		if err != nil {
			continue
		}
		tcpAddrs = append(tcpAddrs, *addr)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tcpAddrs, nil
}

func main() {
	ipList, err := getList("list.txt")
	semaphore := make(chan struct{}, max)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	for _, ip := range ipList {
		wg.Add(2)
		go checker.IsSocks5(ip, &wg, semaphore)
		go checker.IsSocks4(ip, &wg, semaphore)
	}
	wg.Wait()
}

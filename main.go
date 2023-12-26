package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"net"
	"os"
	"proxyChecker/checker"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var outputFile string

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
	var inputFile string
	var maxGoroutines int
	var checkSocks4, checkSocks5 bool
	var dialTimeout, readTimeout time.Duration

	flag.StringVar(&inputFile, "if", "list.txt", "File with list of proxies")
	flag.IntVar(&maxGoroutines, "max", 5000, "Maximum number of goroutines")
	flag.BoolVar(&checkSocks4, "s4", true, "Check for SOCKS4 proxies")
	flag.BoolVar(&checkSocks5, "s5", true, "Check for SOCKS5 proxies")
	flag.DurationVar(&dialTimeout, "dt", 5*time.Second, "Dial timeout duration")
	flag.DurationVar(&readTimeout, "rt", 10*time.Second, "Read timeout duration")
	flag.StringVar(&outputFile, "of", "", "Output file to write results to (defaults to stdout)")
	flag.Parse()

	ipList, err := getList(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	semaphore := make(chan struct{}, maxGoroutines)
	var wg sync.WaitGroup

	var output io.Writer
	var progressBar *pb.ProgressBar
	var counter int64

	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer file.Close()
		output = file

		// Initialize and start the progress bar
		progressBar = pb.StartNew(len(ipList) * (boolToInt(checkSocks4) + boolToInt(checkSocks5)))
	} else {
		output = os.Stdout
	}

	for _, ip := range ipList {
		wg.Add(1)
		go func(ip net.TCPAddr) {
			defer wg.Done()
			if checkSocks4 {
				checker.IsSocks4(ip, semaphore, dialTimeout, readTimeout, output)
				updateProgress(progressBar, &counter)
			}
			if checkSocks5 {
				checker.IsSocks5(ip, semaphore, dialTimeout, readTimeout, output)
				updateProgress(progressBar, &counter)
			}
		}(ip)
	}

	wg.Wait()

	if progressBar != nil {
		progressBar.Finish()
	}
	time.Sleep(1 * time.Second)
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func updateProgress(progressBar *pb.ProgressBar, counter *int64) {
	if progressBar != nil {
		atomic.AddInt64(counter, 1)
		progressBar.SetCurrent(*counter)
	}
}

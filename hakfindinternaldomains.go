package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	concurrencyPtr := flag.Int("t", 8, "Number of threads to utilise. Default is 8.")
	flag.Parse()
	_, cidrA, _ := net.ParseCIDR("10.0.0.0/8")
	_, cidrB, _ := net.ParseCIDR("172.16.0.0/12")
	_, cidrC, _ := net.ParseCIDR("192.168.0.0/16")
	cidrs := [3]net.IPNet{*cidrA, *cidrB, *cidrC}

	work := make(chan string)
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			work <- s.Text()
		}
		close(work)
	}()

	wg := &sync.WaitGroup{}

	for i := 0; i < *concurrencyPtr; i++ {
		wg.Add(1)

		go doWork(work, wg, cidrs)
	}
	wg.Wait()
}

func doWork(work chan string, wg *sync.WaitGroup, cidrs [3]net.IPNet) {
	defer wg.Done()
	for text := range work {
		ip, err := net.LookupIP(text)
		if err != nil {
			log.Println("DNS resolve failed:", err)
		}
		for _, cidr := range cidrs {
			if len(ip) == 0 {
				continue
			}
			if cidr.Contains(ip[0]) {
				fmt.Println(text, ip)
			}
		}
	}
}

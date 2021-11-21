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
	_, cidrD, _ := net.ParseCIDR("0.0.0.0/8")
	_, cidrE, _ := net.ParseCIDR("127.0.0.0/8")
	_, cidrF, _ := net.ParseCIDR("169.254.0.0/16")
	_, cidrG, _ := net.ParseCIDR("192.0.0.0/24")
	_, cidrH, _ := net.ParseCIDR("192.0.2.0/24")
	_, cidrI, _ := net.ParseCIDR("192.88.99.0/24")
	_, cidrJ, _ := net.ParseCIDR("198.18.0.0/15")
	_, cidrK, _ := net.ParseCIDR("198.51.100.0/24")
	_, cidrL, _ := net.ParseCIDR("203.0.113.0/24")
	_, cidrM, _ := net.ParseCIDR("224.0.0.0/4")
	_, cidrN, _ := net.ParseCIDR("240.0.0.0/4")
	_, cidrO, _ := net.ParseCIDR("255.255.255.255/32")
	cidrs := [15]net.IPNet{*cidrA, *cidrB, *cidrC, *cidrD, *cidrE, *cidrF, *cidrG, *cidrH, *cidrI, *cidrJ, *cidrK, *cidrL, *cidrM, *cidrN, *cidrO}

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

func doWork(work chan string, wg *sync.WaitGroup, cidrs [15]net.IPNet) {
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

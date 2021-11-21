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
	cidr_strings := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"0.0.0.0/8",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"192.88.99.0/24",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"224.0.0.0/4",
		"240.0.0.0/4",
		"255.255.255.255/32",
	}
	cidrs := make([]net.IPNet, 15)
	for _, cidr_str := range cidr_strings {
		_, cidr, _ := net.ParseCIDR(cidr_str)
		cidrs = append(cidrs, *cidr)
	}

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

func doWork(work chan string, wg *sync.WaitGroup, cidrs []net.IPNet) {
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

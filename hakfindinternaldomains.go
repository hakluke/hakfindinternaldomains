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
	// Taken from https://datatracker.ietf.org/doc/html/rfc5735
	cidr_strings := []string{
		"0.0.0.0/8",          // "This" Network RFC 1122, Section 3.2.1.3
		"10.0.0.0/8",         // Private-Use Networks RFC 1918
		"127.0.0.0/8",        // Loopback RFC 1122, Section 3.2.1.3
		"169.254.0.0/16",     // Link Local RFC 3927
		"172.16.0.0/12",      // Private-Use Networks RFC 1918
		"192.0.0.0/24",       // IETF Protocol Assignments RFC 5736
		"192.0.2.0/24",       // TEST-NET-1 RFC 5737
		"192.88.99.0/24",     // 6to4 Relay Anycast RFC 3068
		"192.168.0.0/16",     // Private-Use Networks RFC 1918
		"198.18.0.0/15",      // Network Interconnect
		"198.51.100.0/24",    // TEST-NET-2 RFC 5737
		"203.0.113.0/24",     // TEST-NET-3 RFC 5737
		"224.0.0.0/4",        // Multicast RFC 3171
		"240.0.0.0/4",        // Reserved for Future Use RFC 1112, Section 4
		"255.255.255.255/32", // Limited Broadcast RFC 919, Section 7
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

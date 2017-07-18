package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sparrc/go-ping"
)

var (
	host string
)

func init() {
	flag.StringVar(&host, "host", "", "the hostname to ping")
}

func main() {
	flag.Parse()

	if host == "" {
		log.Print("a hostname is required")
		flag.Usage()
		os.Exit(1)
	}
	pinger, err := ping.NewPinger(host)
	if err != nil {
		log.Fatalf("error creating the pinger: %s", err)
	}

	// tell pinger that it is privileged.
	// NOTE: You must run `setcap cap_net_raw=+ep pocket-loss-monitor`
	pinger.SetPrivileged(true)

	var lastPing int
	var lostPackets []int
	pinger.OnRecv = func(pkt *ping.Packet) {
		for i := lastPing + 1; i < pkt.Seq; i++ {
			lostPackets = append(lostPackets, i)
		}
		lastPing = pkt.Seq
		if pkt.Seq > 0 && pkt.Seq%10 == 0 {
			fmt.Printf("lost packets: %v\n", lostPackets)
		}
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("lost packets: %v\n", lostPackets)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()
}

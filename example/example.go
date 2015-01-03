// example
package main

import (
	"flag"
	"fmt"
	"github.com/subh007/goodl/go_ping_sweep"
	"os"
)

func main() {
	host := flag.String("host", "", "ip address/ hostname to analyse the ping.")
	flag.Parse()

	if host == nil {
		fmt.Println("usage: ./example --host <host name/ip>")
		os.Exit(-1)
	}

	if go_ping_sweep.IsAdmin() {
		t, err := go_ping_sweep.PingAnalyse(*host, 10)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}
		t.CreateTable()
	}
}

This project analyse the RTT time for a ipv4 host or ipv6 host using the ICMP echo request/ reply message.

It is try to port the [ping_sweep] (https://github.com/Who8MyLunch/Ping_Sweep) project in go language which is written in the Python language.

Example
=======
[example.go] (https://github.com/subh007/go_ping_sweep/blob/master/example/example.go)
```go
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
```

After building and installing (`go build && go install` in example folder) the example code, run the command with the root priviledge :

```shell
$ sudo example --host google.com
TimePing | DataSize | PacketSize | status |
79.16993ms | - | - | 1 |
65.865353ms | - | - | 1 |
22.95475ms | - | - | 1 |
25.013145ms | - | - | 1 |
9.348976ms | - | - | 1 |
21.015361ms | - | - | 1 |
31.260762ms | - | - | 1 |
33.435962ms | - | - | 1 |
22.23161ms | - | - | 1 |
28.16272ms | - | - | 1 |
```
TimePing : RTT(round trip time), time taken to receive the ICMP echo reply after sending ICMP echo request message.

DataSize : payload size of ICMP packet.

PacketSize : packet size of including ICMP header.

status : PING PASS/FAIL.

Feature to be added:
-------------------
- Non-blocking echo request-response (using go-routines and workgroup).
- ipv6 ping support.
- ping analysis based on the packet size.
- beautify the table alignment with the table header (I have plan to move the [table.go] (https://github.com/subh007/go_ping_sweep/blob/master/table.go) as standalone project).
- bugfixes

Please feel free to send comment/ feedback / suggestion. 

Reference
---------
https://github.com/Who8MyLunch/Ping_Sweep

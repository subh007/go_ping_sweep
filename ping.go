package go_ping_sweep

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// structure to hold the ping result.
type Result struct {
	TimePing   int64 // rtt time in ns
	DataSize   int   // data size in icmp packet
	PacketSize int   // packet size of icmp
	Status     bool  // status for ping pass/fail
}

// function sends the single icmp packet and respnse back with the
// Result.
func singlePing(host string, conn net.Conn) (*Result, error) {
	xid, xseq := os.Getpid()&0xffff, 1

	// create icmp packet
	icmp := icmpMessage{
		Type: 8,
		Code: 0,
		Body: &icmpMessageBody{
			ID: xid, Seq: xseq,
			Data: []byte("Go Go packet"),
		},
	}

	icmp_byte, err := icmp.Marshal()
	if err != nil {
		return nil, err
	}

	send_time := time.Now()
	_, err = conn.Write(icmp_byte)

	if err != nil {
		fmt.Println("err: " + err.Error())
		return nil, err
	}

	// capture the ping response message
	// NOTE: extra 20 bytes is asssigned becuaes the
	// response packet will carry the IP header as well.
	rb := make([]byte, 20+len(icmp_byte))

	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	if _, err = conn.Read(rb); err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	// first check that packet is received with correct
	// (xid, xseq). If not then it is not the correct response
	// and wait for the some timeout period.
	icmp_rep, err := parseICMPMessage(rb[20:])
	if err != nil {
		fmt.Println("error generated")
		return nil, err
	}

	var r Result

	if icmp_rep.Body.ID == xid && icmp_rep.Body.Seq == xseq {
		rcvd_time := time.Now()
		r.TimePing = rcvd_time.Sub(send_time).Nanoseconds()
	} else {
		fmt.Println("not matching")
	}
	r.DataSize = 0
	r.PacketSize = 0
	r.Status = true

	return &r, nil
}

// setup the connection for sending the ICMP packet.
func setupConnection(conn_type, host string) (net.Conn, error) {
	addrs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	host = addrs[0].String()
	conn, err := net.Dial(conn_type, host)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// funtion analyze the ping time and return in table format.
func PingAnalyse(host string, pkt_count int) (*Table, error) {
	t := new(Table)
	t.SetHeader("TimePing (ns)", "DataSize", "PacketSize", "status", "mean")

	conn, err := setupConnection("ip4:icmp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	results := make([]Result, pkt_count)

	for i := 0; i < pkt_count; i++ {
		res, err := singlePing(host, conn)
		results[i] = *res

		if err != nil {
			t.AddData("0", "0", "0", "-1")
		} else {
			t.AddData(
				strconv.FormatInt(res.TimePing, 10),
				"-",
				"-",
				"1",
			)
		}
	}

	fmt.Println(average(results))
	return t, nil
}

/*
// testing the icmp sniffer
func SniffICMP() {
	conn, err := setupConnection("ip4:icmp", "google.com")
	if err != nil {
		fmt.Println(err.Error())
	}
	rb := make([]byte, 100)

	for {
		if _, err = conn.Read(rb); err != nil {
			fmt.Print(err.Error())
		}
		fmt.Println(string(rb))
	}
}
*/

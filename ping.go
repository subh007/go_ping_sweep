package go_ping_sweep

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type icmpMessageBody struct {
	ID   int    // identifier
	Seq  int    // sequence number
	Data []byte // data
}

// icmp packet
type icmpMessage struct {
	Type     int              // type
	Code     int              // code
	Checksum int              // checksum
	Body     *icmpMessageBody // body
}

func (m *icmpMessage) Marshal() ([]byte, error) {
	b := []byte{byte(m.Type), byte(m.Code), 0, 0}
	if m.Body != nil && m.Body.Len() != 0 {
		mb, err := m.Body.Marshal()
		if err != nil {
			return nil, err
		}
		b = append(b, mb...)
	}
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	// Place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero.
	b[2] ^= byte(^s)
	b[3] ^= byte(^s >> 8)
	return b, nil
}

// Marshal returns the binary enconding of the ICMP echo request or
// reply message body p.
func (p *icmpMessageBody) Marshal() ([]byte, error) {
	b := make([]byte, 4+len(p.Data))
	b[0], b[1] = byte(p.ID>>8), byte(p.ID)
	b[2], b[3] = byte(p.Seq>>8), byte(p.Seq)
	copy(b[4:], p.Data)
	return b, nil
}

func (p *icmpMessageBody) Len() int {
	return 4 + len(p.Data)
}

func parseICMPMessageBody(b []byte) (*icmpMessageBody, error) {
	p := &icmpMessageBody{
		ID:  int(b[0]<<8) | int(b[1]),
		Seq: int(b[2]<<8) | int(b[3]),
	}

	p.Data = make([]byte, len(b)-4)
	copy(p.Data, b[4:])
	return p, nil
}

func parseICMPMessage(b []byte) (*icmpMessage, error) {
	m := &icmpMessage{
		Type:     int(b[0]),
		Code:     int(b[1]),
		Checksum: int(b[2]<<8) | int(b[3]),
	}

	var err error
	m.Body, err = parseICMPMessageBody(b[4:])
	if err != nil {
		fmt.Println("message can't be parsed")
		return nil, err
	}
	return m, nil
}

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
	rb := make([]byte, 40+len(icmp_byte))

	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	if _, err = conn.Read(rb); err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	// TODO: first check that packet is received with correct
	// (xid, xseq). If not then it is not the correct response
	// and wait for the some timeout period. It is important
	// because it might possible we are receiving ICMP response
	// for other process. We have to confirm that would it be possible
	// to sniff all the ICMP packets here.
	rcvd_time := time.Now()

	diff := rcvd_time.Sub(send_time).Nanoseconds()
	_, err = parseICMPMessage(rb)
	if err != nil {
		fmt.Println("err: " + err.Error())
		return nil, err
	}

	var r Result
	r.TimePing = diff
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
	conn, err := net.Dial("ip4:icmp", host)
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

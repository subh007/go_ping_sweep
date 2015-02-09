package go_ping_sweep

import (
	"fmt"
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
		ID:  (int(b[0]) << 8) | int(b[1]),
		Seq: (int(b[2]) << 8) | int(b[3]),
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

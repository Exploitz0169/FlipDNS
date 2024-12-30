package parser

import (
	"encoding/binary"
	"errors"
)

type DNSQuestion struct {
	// a domain name represented as a sequence of labels, where
	// each label consists of a length octet followed by that
	// number of octets.  The domain name terminates with the
	// zero length octet for the null label of the root.  Note
	// that this field may be an odd number of octets; no
	// padding is used.
	QNAME []byte
	// a two octet code which specifies the type of the query.
	// The values for this field include all codes valid for a
	// TYPE field, together with some more general codes which
	// can match more than one type of RR.
	QTYPE uint16
	// a two octet code that specifies the class of the query.
	// For example, the QCLASS field is IN for the Internet
	QCLASS uint16
}

var (
	ErrInvalidQuestionLength = errors.New("invalid DNS question length")
	ErrBufferTooShort        = errors.New("buffer too short")
)

func ParseDNSQuestions(buf []byte, qdcount uint16) ([]*DNSQuestion, error) {

	if qdcount == 0 || len(buf) == 0 {
		return nil, ErrInvalidQuestionLength
	}

	questions := make([]*DNSQuestion, qdcount)
	BE := binary.BigEndian.Uint16

	offset := 0
	for i := 0; i < int(qdcount); i++ {
		qname, n := parseQNAME(buf[offset:])
		questions[i] = &DNSQuestion{
			QNAME:  qname,
			QTYPE:  BE(buf[offset+n : offset+n+2]),
			QCLASS: BE(buf[offset+n+2 : offset+n+4]),
		}
		offset += n + 4
	}

	return questions, nil
}

// Return the QNAME and the number of bytes read
func parseQNAME(buf []byte) ([]byte, int) {
	n := 0
	for i := 0; i < len(buf); i++ {
		// The domain name terminates with the
		// zero length octet for the null label of the root.
		if buf[i] == 0x00 {
			n = i + 1
			break
		}
	}
	return buf[:n], n
}
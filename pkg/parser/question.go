package parser

import (
	"bytes"
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
	// Domain name is parsed from QNAME, it is easiest to just
	// return it alongside the DNS Question
	DOMAIN string
	// a two octet code which specifies the type of the query.
	// The values for this field include all codes valid for a
	// TYPE field, together with some more general codes which
	// can match more than one type of RR.
	QTYPE uint16
	// a two octet code that specifies the class of the query.
	// For example, the QCLASS field is IN for the Internet
	QCLASS uint16
}

func (q *DNSQuestion) Serialize() ([]byte, error) {
	buf := bytes.Buffer{}

	if _, err := buf.Write(q.QNAME); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, q.QTYPE); err != nil {
		return nil, err
	}

	if err := binary.Write(&buf, binary.BigEndian, q.QCLASS); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var (
	ErrInvalidQuestionLength = errors.New("invalid DNS question length")
	ErrBufferTooShort        = errors.New("buffer too short")
	ErrInvalidQNAME          = errors.New("invalid QNAME. Check bytes and length octet")
	ErrEmptyQNAME            = errors.New("invalid QNAME. No domain name was parsed from")
)

// Expects buf to not include DNS Header
func ParseDNSQuestions(buf []byte, qdcount uint16) ([]*DNSQuestion, error) {

	if qdcount == 0 || len(buf) == 0 {
		return nil, ErrInvalidQuestionLength
	}

	questions := make([]*DNSQuestion, qdcount)
	BE := binary.BigEndian.Uint16

	offset := 0
	for i := 0; i < int(qdcount); i++ {
		qname, n := parseQNAME(buf[offset:])
		domain, err := parseDomainFromQNAME(qname)
		if err != nil {
			return nil, err
		}
		questions[i] = &DNSQuestion{
			QNAME:  qname,
			QTYPE:  BE(buf[offset+n : offset+n+2]),
			QCLASS: BE(buf[offset+n+2 : offset+n+4]),
			DOMAIN: domain,
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

func parseDomainFromQNAME(qname []byte) (string, error) {
	// Expecting at least 3. The first length octect, the label, and the terminator
	if len(qname) < 3 {
		return "", ErrBufferTooShort
	}
	domain := ""
	i := 0
	for i < len(qname) {
		// QNAME ends in zero length octet
		if qname[i] == 0x00 {
			break
		}
		// Sequence of labels where each label consists of a length octet
		// followed by that number of octets
		labelLength := int(qname[i])
		if labelLength < 0 || i+1+labelLength >= len(qname) {
			return "", ErrInvalidQNAME
		}
		domain += string(qname[i+1:i+1+labelLength]) + "."
		i += labelLength + 1
	}
	if len(domain) < 1 {
		return "", ErrEmptyQNAME
	}
	return domain[:len(domain)-1], nil
}

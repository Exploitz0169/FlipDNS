package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

type DNSResourceRecord struct {

	// a domain name to which this resource record pertains.
	NAME []byte
	// two octets containing one of the RR type codes.  This
	// field specifies the meaning of the data in the RDATA
	// field.
	TYPE uint16
	// two octets which specify the class of the data in the
	// RDATA field.
	CLASS uint16
	// a 32 bit unsigned integer that specifies the time
	// interval (in seconds) that the resource record may be
	// cached before it should be discarded.  Zero values are
	// interpreted to mean that the RR can only be used for the
	// transaction in progress, and should not be cached.
	TTL uint32
	// an unsigned 16 bit integer that specifies the length in
	// octets of the RDATA field.
	RDLENGTH uint16
	// a variable length string of octets that describes the
	// resource.  The format of this information varies
	// according to the TYPE and CLASS of the resource record.
	// For example, the if the TYPE is A and the CLASS is IN,
	// the RDATA field is a 4 octet ARPA Internet address.
	RDATA []byte
}

type DNSClass = uint16

const (
	// the Internet
	ClassIN DNSClass = 1
	// the CSNET class (Obsolete - used only for examples in
	// some obsolete RFCs)
	ClassCS DNSClass = 2
	// the CHAOS class
	ClassCH DNSClass = 3
	// Hesiod [Dyer 87]
	ClassHS DNSClass = 4
	// any class
	ClassANY DNSClass = 255
)

type DNSType = uint16

const (
	// a host address
	TypeA DNSType = 1
	// an authoritative name server
	TypeNS DNSType = 2
	// a mail destination (Obsolete - use MX)
	TypeMD DNSType = 3
	// a mail forwarder (Obsolete - use MX)
	TypeMF DNSType = 4
	// the canonical name for an alias
	TypeCNAME DNSType = 5
	// marks the start of a zone of authority
	TypeSOA DNSType = 6
	// a mailbox domain name (EXPERIMENTAL)
	TypeMB DNSType = 7
	// a mail group member (EXPERIMENTAL)
	TypeMG DNSType = 8
	// a mail rename domain name (EXPERIMENTAL)
	TypeMR DNSType = 9
	// a null RR (EXPERIMENTAL)
	TypeNULL DNSType = 10
	// a well known service description
	TypeWKS DNSType = 11
	// a domain name pointer
	TypePTR DNSType = 12
	// host information
	TypeHINFO DNSType = 13
	// mailbox or mail list information
	TypeMINFO DNSType = 14
	// mail exchange
	TypeMX DNSType = 15
	// text strings
	TypeTXT DNSType = 16
	// IPv6 address
	TypeAAAA DNSType = 28
	// service locator
	TypeSRV DNSType = 33
	// naming authority pointer
	TypeNAPTR DNSType = 35
	// option
	TypeOPT DNSType = 41
	// transfer of an entire zone
	TypeAXFR DNSType = 252
	// all records
	TypeALL DNSType = 255
)

func (rr *DNSResourceRecord) Serialize() ([]byte, error) {
	buf := bytes.Buffer{}

	if _, err := buf.Write(rr.NAME); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, rr.TYPE); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, rr.CLASS); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, rr.TTL); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, rr.RDLENGTH); err != nil {
		return nil, err
	}
	if _, err := buf.Write(rr.RDATA); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var (
	ErrInvalidIP = errors.New("invalid IP to parse to RDATA")
)

func CreateDNSAAnswer(owner []byte, ipv4 string, ttl uint32) (*DNSResourceRecord, error) {

	if err := ValidateDNSName(owner); err != nil {
		return nil, err
	}

	ip, err := getIpv4AsBytes(ipv4)
	if err != nil {
		return nil, err
	}

	return &DNSResourceRecord{
		NAME:     owner,
		TYPE:     TypeA,
		CLASS:    ClassIN,
		TTL:      ttl,
		RDLENGTH: 4,
		RDATA:    ip,
	}, nil

}

func getIpv4AsBytes(domain string) ([]byte, error) {
	sections := strings.Split(domain, ".")
	if len(sections) != 4 {
		return nil, ErrInvalidIP
	}

	bytes := make([]byte, 4)
	for i, section := range sections {
		section_int, err := strconv.Atoi(section)
		if err != nil {
			return nil, ErrInvalidIP
		}
		bytes[i] = byte(section_int)
	}

	return bytes, nil
}

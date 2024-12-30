package parser

import (
	"encoding/binary"
	"errors"
)

type DNSHeaderFlags struct {
	// A one bit field that specifies whether this message is a
	// query (0), or a response (1).
	QR uint8
	// A four bit field that specifies kind of query in this
	// message.  This value is set by the originator of a query
	// and copied into the response.  The values are:

	// 0               a standard query (QUERY)

	// 1               an inverse query (IQUERY)

	// 2               a server status request (STATUS)

	// 3-15            reserved for future use
	OPCODE uint8
	// Authoritative Answer - this bit is valid in responses,
	// and specifies that the responding name server is an
	// authority for the domain name in question section.

	// Note that the contents of the answer section may have
	// multiple owner names because of aliases.  The AA bit
	// corresponds to the name which matches the query name, or
	// the first owner name in the answer section.
	AA uint8
	// TrunCation - specifies that this message was truncated
	// due to length greater than that permitted on the
	// transmission channel.
	TC uint8
	// Recursion Desired - this bit may be set in a query and
	// is copied into the response.  If RD is set, it directs
	// the name server to pursue the query recursively.
	// Recursive query support is optional.
	RD uint8
	// Recursion Available - this be is set or cleared in a
	// response, and denotes whether recursive query support is
	// available in the name server.
	RA uint8
	// Reserved for future use.  Must be zero in all queries
	// and responses.
	Z uint8
	// Response code - this 4 bit field is set as part of
	// responses.  The values have the following
	// interpretation:

	// 0               No error condition

	// 1               Format error - The name server was
	// 				unable to interpret the query.

	// 2               Server failure - The name server was
	// 				unable to process this query due to a
	// 				problem with the name server.

	// 3               Name Error - Meaningful only for
	// 				responses from an authoritative name
	// 				server, this code signifies that the
	// 				domain name referenced in the query does
	// 				not exist.

	// 4               Not Implemented - The name server does
	// 				not support the requested kind of query.

	// 5               Refused - The name server refuses to
	// 				perform the specified operation for
	// 				policy reasons.  For example, a name
	// 				server may not wish to provide the
	// 				information to the particular requester,
	// 				or a name server may not wish to perform
	// 				a particular operation (e.g., zone
	// 				transfer) for particular data.
	// 6-15            Reserved for future use.
	RCODE uint8
}

type DNSHeader struct {
	// A 16 bit identifier assigned by the program that
	// generates any kind of query.  This identifier is copied
	// the corresponding reply and can be used by the requester
	// to match up replies to outstanding queries.
	ID    uint16
	Flags *DNSHeaderFlags
	// an unsigned 16 bit integer specifying the number of
	// entries in the question section.
	QDCOUNT uint16
	// an unsigned 16 bit integer specifying the number of
	// resource records in the answer section.
	ANCOUNT uint16
	// an unsigned 16 bit integer specifying the number of name
	// server resource records in the authority records
	// section.
	NSCOUNT uint16
	// an unsigned 16 bit integer specifying the number of
	// resource records in the additional records section.
	ARCOUNT uint16
}

var (
	ErrInvalidHeaderLength = errors.New("invalid DNS header length: must be 12 bytes")
	ErrInvalidFlagsLength  = errors.New("invalid DNS flags length: must be 2 bytes")
)

func ParseDNSHeader(buf []byte) (*DNSHeader, error) {

	if len(buf) != 12 {
		return nil, ErrInvalidHeaderLength
	}

	BE := binary.BigEndian.Uint16
	flags, err := ParseDNSHeaderFlags(buf[2:4])
	if err != nil {
		return nil, err
	}

	header := &DNSHeader{
		ID:      BE(buf[0:2]),
		Flags:   flags,
		QDCOUNT: BE(buf[4:6]),
		ANCOUNT: BE(buf[6:8]),
		NSCOUNT: BE(buf[8:10]),
		ARCOUNT: BE(buf[10:12]),
	}

	return header, nil
}

func ParseDNSHeaderFlags(buf []byte) (*DNSHeaderFlags, error) {

	if len(buf) != 2 {
		return nil, ErrInvalidFlagsLength
	}

	return &DNSHeaderFlags{
		QR:     buf[0] >> 7,
		OPCODE: buf[0] & 0b01111000 >> 3,
		AA:     buf[0] & 0b00000100 >> 2,
		TC:     buf[0] & 0b00000010 >> 1,
		RD:     buf[0] & 0b00000001,
		RA:     buf[1] >> 7,
		Z:      buf[1] & 0b01110000 >> 4,
		RCODE:  buf[1] & 0b00001111,
	}, nil
}

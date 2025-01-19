package dns

func ParseDNSQuery(buf []byte) (*DNSHeader, []*DNSQuestion, error) {
	header, err := ParseDNSHeader(buf[:12])
	if err != nil {
		return nil, nil, err
	}

	questions, err := ParseDNSQuestions(buf[12:], header.QDCOUNT)
	if err != nil {
		return nil, nil, err
	}

	return header, questions, nil
}

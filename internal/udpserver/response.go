package udpserver

import "github.com/exploitz0169/flipdns/pkg/dns"

func (s *UdpServer) buildResponse(
	header *dns.DNSHeader,
	questions []*dns.DNSQuestion,
	answers []*dns.DNSResourceRecord,
	rcode uint8,
) ([]byte, error) {
	responseHeader, err := dns.CreateDNSAnswerHeader(header, uint16(len(answers)), 0, 0, true, rcode)
	if err != nil {
		return nil, err
	}

	items := make([]Serializable, 0)
	items = append(items, responseHeader)
	for _, question := range questions {
		items = append(items, question)
	}
	for _, answer := range answers {
		items = append(items, answer)
	}

	return SerializeItems(items)
}

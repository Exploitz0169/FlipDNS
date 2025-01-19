package udpserver

import (
	"log/slog"
	"net"

	"github.com/exploitz0169/flipdns/pkg/dns"
)

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

func (s *UdpServer) sendErrorResponse(
	conn net.PacketConn,
	addr net.Addr,
	header *dns.DNSHeader,
	questions []*dns.DNSQuestion,
	rcode uint8,
) error {
	response, err := s.buildResponse(header, questions, nil, rcode)
	if err != nil {
		s.app.Logger.Warn("Failed to build response",
			slog.String("error", err.Error()),
		)
		return err
	}

	if err = s.sendResponse(conn, addr, response); err != nil {
		s.app.Logger.Warn("Failed to send response",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (s *UdpServer) sendResponse(conn net.PacketConn, addr net.Addr, packet []byte) error {
	_, err := conn.WriteTo(packet, addr)
	if err != nil {
		return err
	}

	return nil
}

package udpserver

import (
	"log/slog"
	"net"

	"github.com/exploitz0169/flipdns/pkg/dns"
)

func (s *UdpServer) handlePacket(buf []byte, addr net.Addr, conn net.PacketConn) {
	header, questions, err := dns.ParseDNSQuery(buf)
	if err != nil {
		s.app.Logger.Warn("Failed to parse DNS query",
			slog.String("error", err.Error()),
		)
		s.sendErrorResponse(conn, addr, header, questions, dns.RCodeFormatError)
		return
	}

	answers, err := s.resolveQuestions(questions)
	if err != nil {
		s.app.Logger.Warn("Failed to resolve questions",
			slog.String("error", err.Error()),
		)

		if err == ErrRecordNotFound {
			s.sendErrorResponse(conn, addr, header, questions, dns.RCodeNameError)
		} else {
			s.sendErrorResponse(conn, addr, header, questions, dns.RCodeServerFail)
		}

		return
	}

	packet, err := s.buildResponse(header, questions, answers, 0)
	if err != nil {
		s.app.Logger.Warn("Failed to build response",
			slog.String("error", err.Error()),
		)
		s.sendErrorResponse(conn, addr, header, questions, dns.RCodeServerFail)
		return
	}

	if err = s.sendResponse(conn, addr, packet); err != nil {
		s.app.Logger.Warn("Failed to send response",
			slog.String("error", err.Error()),
		)
	}
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

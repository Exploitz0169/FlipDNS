package udpserver

import (
	"log/slog"
	"net"

	"github.com/exploitz0169/flipdns/pkg/dns"
)

func (s *UdpServer) handlePacket(buf []byte, addr net.Addr) {
	header, questions, err := dns.ParseDNSQuery(buf)
	if err != nil {
		s.app.Logger.Warn("Failed to parse DNS query",
			slog.String("error", err.Error()),
		)
		s.sendErrorResponse(addr, header, questions, dns.RCodeFormatError)
		return
	}

	answers, err := s.resolveQuestions(questions)
	if err != nil {
		s.app.Logger.Warn("Failed to resolve questions",
			slog.String("error", err.Error()),
		)

		if err == ErrRecordNotFound {
			s.sendErrorResponse(addr, header, questions, dns.RCodeNameError)
		} else {
			s.sendErrorResponse(addr, header, questions, dns.RCodeServerFail)
		}

		return
	}

	packet, err := s.buildResponse(header, questions, answers, 0)
	if err != nil {
		s.app.Logger.Warn("Failed to build response",
			slog.String("error", err.Error()),
		)
		s.sendErrorResponse(addr, header, questions, dns.RCodeServerFail)
		return
	}

	if err = s.sendResponse(addr, packet); err != nil {
		s.app.Logger.Warn("Failed to send response",
			slog.String("error", err.Error()),
		)
	}
}

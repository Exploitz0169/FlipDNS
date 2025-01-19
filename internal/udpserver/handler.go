package udpserver

import (
	"errors"
	"log/slog"
	"net"

	"github.com/exploitz0169/flipdns/pkg/dns"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

func (s *UdpServer) handlePacket(buf []byte, addr net.Addr, conn net.PacketConn) {

	header, questions, err := dns.ParseDNSQuery(buf)
	if err != nil {
		s.app.Logger.Warn("Failed to parse DNS query", slog.String("error", err.Error()))
		s.sendErrorResponse(conn, addr, header, questions, dns.RCodeFormatError)
		return
	}

	answers, err := s.resolveQuestions(questions)
	if err != nil {
		s.app.Logger.Warn("Failed to resolve questions", slog.String("error", err.Error()))

		if err == ErrRecordNotFound {
			s.sendErrorResponse(conn, addr, header, questions, dns.RCodeNameError)
		} else {
			s.sendErrorResponse(conn, addr, header, questions, dns.RCodeServerFail)
		}

		return
	}

	packet, err := s.buildResponse(header, questions, answers, 0)
	if err != nil {
		s.app.Logger.Warn("Failed to build response", slog.String("error", err.Error()))
		s.sendErrorResponse(conn, addr, header, questions, dns.RCodeServerFail)
		return
	}

	if err = s.sendResponse(conn, addr, packet); err != nil {
		s.app.Logger.Warn("Failed to send response", slog.String("error", err.Error()))
	}
}

func (s *UdpServer) resolveQuestions(questions []*dns.DNSQuestion) ([]*dns.DNSResourceRecord, error) {
	answers := make([]*dns.DNSResourceRecord, 0, len(questions))

	for _, question := range questions {

		s.app.Logger.Info("Resolving question", slog.String("domain", string(question.DOMAIN)))

		record, ok := s.app.Db.GetRecord(question.DOMAIN)
		if !ok {
			s.app.Logger.Warn("Record not found", slog.String("domain", string(question.DOMAIN)))
			return nil, ErrRecordNotFound
		}

		answer, err := dns.CreateDNSAAnswer(question.QNAME, record.IPv4, record.TTL)
		if err != nil {
			s.app.Logger.Warn("Error creating DNS A answer", slog.String("error", err.Error()))
			return nil, err
		}

		answers = append(answers, answer)
	}

	return answers, nil
}

func (s *UdpServer) buildResponse(header *dns.DNSHeader, questions []*dns.DNSQuestion, answers []*dns.DNSResourceRecord, rcode uint8) ([]byte, error) {
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

func (s *UdpServer) sendErrorResponse(conn net.PacketConn, addr net.Addr, header *dns.DNSHeader, questions []*dns.DNSQuestion, rcode uint8) error {

	response, err := s.buildResponse(header, questions, nil, rcode)
	if err != nil {
		s.app.Logger.Warn("Failed to build response", slog.String("error", err.Error()))
		return err
	}

	if err = s.sendResponse(conn, addr, response); err != nil {
		s.app.Logger.Warn("Failed to send response", slog.String("error", err.Error()))
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

package udpserver

import (
	"errors"
	"log/slog"

	"github.com/exploitz0169/flipdns/pkg/dns"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

func (s *UdpServer) resolveQuestions(questions []*dns.DNSQuestion) ([]*dns.DNSResourceRecord, error) {
	answers := make([]*dns.DNSResourceRecord, 0, len(questions))

	for _, question := range questions {
		s.app.Logger.Info("Resolving question",
			slog.String("domain", string(question.DOMAIN)),
		)

		record, ok := s.app.Db.GetRecord(question.DOMAIN)
		if !ok {
			s.app.Logger.Warn("Record not found",
				slog.String("domain", string(question.DOMAIN)),
			)
			return nil, ErrRecordNotFound
		}

		answer, err := dns.CreateDNSAAnswer(question.QNAME, record.IPv4, record.TTL)
		if err != nil {
			s.app.Logger.Warn("Error creating DNS A answer",
				slog.String("error", err.Error()),
			)
			return nil, err
		}

		answers = append(answers, answer)
	}

	return answers, nil
}

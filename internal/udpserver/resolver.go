package udpserver

import (
	"context"
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

		record, error := s.app.Db.GetRecordByDomainName(context.Background(), question.DOMAIN)
		if error != nil {
			s.app.Logger.Warn("Record not found",
				slog.String("domain", string(question.DOMAIN)),
			)
			return nil, ErrRecordNotFound
		}

		answer, err := dns.CreateDNSAAnswer(question.QNAME, record.RecordData, uint32(record.Ttl))
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

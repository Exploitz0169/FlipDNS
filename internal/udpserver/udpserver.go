package udpserver

import (
	"bytes"
	"log/slog"
	"net"

	"github.com/exploitz0169/flipdns/pkg/parser"
)

// Temp just to test the parser
func Run(conn net.PacketConn) {
	buf := make([]byte, 512)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			slog.Warn("Failed to read packet",
				slog.Int("bytes", n),
				slog.String("addr", addr.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		slog.Info("Received packet",
			slog.Int("bytes", n),
			slog.String("addr", addr.String()),
		)

		header, err := parser.ParseDNSHeader(buf[:12])
		if err != nil {
			slog.Warn("Error parsing DNS header", slog.String("error", err.Error()))
			continue
		}

		questions, err := parser.ParseDNSQuestions(buf[12:], header.QDCOUNT)
		if err != nil {
			slog.Warn("Error parsing DNS questions", slog.String("error", err.Error()))
			continue
		}

		question := questions[0]

		answer, err := parser.CreateDNSAAnswer(question.QNAME, "192.168.2.143", 300)
		if err != nil {
			slog.Warn("Error creating DNS A answer", slog.String("error", err.Error()))
			continue
		}

		responseHeader, err := parser.CreateDNSAnswerHeader(header, 1, 0, 0)
		if err != nil {
			slog.Warn("Error creating DNS answer header", slog.String("error", err.Error()))
			continue
		}

		items := make([]Serializable, 0, 3)
		items = append(items, responseHeader, question, answer)

		packet, err := SerializeItems(items)
		if err != nil {
			slog.Warn("Error serializing items", slog.String("error", err.Error()))
			continue
		}

		_, err = conn.WriteTo(packet, addr)
		if err != nil {
			slog.Warn("Error writing packet", slog.String("error", err.Error()))
			continue
		}

	}
}

type Serializable interface {
	Serialize() ([]byte, error)
}

func SerializeItems(items []Serializable) ([]byte, error) {

	buf := bytes.Buffer{}

	for _, item := range items {
		serialized, err := item.Serialize()
		if err != nil {
			return nil, err
		}
		buf.Write(serialized)
	}

	return buf.Bytes(), nil
}

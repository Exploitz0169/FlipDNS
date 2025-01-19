package udpserver

import (
	"log/slog"
	"net"

	"github.com/exploitz0169/flipdns/internal/app"
)

type UdpServer struct {
	app *app.App
}

func NewUdpServer(app *app.App) *UdpServer {
	return &UdpServer{
		app: app,
	}
}

// Temp just to test the parser
func (s *UdpServer) Run(conn net.PacketConn) {
	buf := make([]byte, 512)

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			s.app.Logger.Warn("Failed to read packet",
				slog.Int("bytes", n),
				slog.String("addr", addr.String()),
				slog.String("error", err.Error()),
			)
			continue
		}

		s.app.Logger.Info("Received packet",
			slog.Int("bytes", n),
			slog.String("addr", addr.String()),
		)

		go s.handlePacket(buf[:n], addr, conn)
	}
}

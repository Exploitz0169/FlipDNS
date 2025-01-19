package udpserver

import (
	"log/slog"
	"net"

	"github.com/exploitz0169/flipdns/internal/app"
)

type UdpServer struct {
	app  *app.App
	conn net.PacketConn
}

func NewUdpServer(app *app.App, conn net.PacketConn) *UdpServer {
	return &UdpServer{
		app:  app,
		conn: conn,
	}
}

func (s *UdpServer) Run() {
	buf := make([]byte, 512)

	for {
		n, addr, err := s.conn.ReadFrom(buf)
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

		go s.handlePacket(buf[:n], addr)
	}
}

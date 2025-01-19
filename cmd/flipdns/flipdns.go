package main

import (
	"log/slog"
	"net"

	"github.com/exploitz0169/flipdns/internal/logger"
	"github.com/exploitz0169/flipdns/internal/udpserver"
)

func main() {

	logger.InitLogger()

	addr := ":53"

	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer conn.Close()

	slog.Info("Started UDP server on addr " + addr)
	udpserver.Run(conn)

}

package main

import (
	"net"

	"github.com/exploitz0169/flipdns/internal/app"
	"github.com/exploitz0169/flipdns/internal/database"
	"github.com/exploitz0169/flipdns/internal/logger"
	"github.com/exploitz0169/flipdns/internal/udpserver"
)

func main() {

	logger := logger.NewLogger()

	app := &app.App{
		Db:     database.NewDatabase(),
		Logger: logger,
	}

	addr := ":53"

	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer conn.Close()

	logger.Info("Started UDP server on addr " + addr)

	server := udpserver.NewUdpServer(app, conn)
	server.Run()

}

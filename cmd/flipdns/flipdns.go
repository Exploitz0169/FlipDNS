package main

import (
	"context"
	"net"

	"github.com/exploitz0169/flipdns/internal/app"
	"github.com/exploitz0169/flipdns/internal/logger"
	"github.com/exploitz0169/flipdns/internal/repository"
	"github.com/exploitz0169/flipdns/internal/udpserver"
	"github.com/jackc/pgx/v5"
)

func main() {

	ctx := context.Background()

	logger := logger.NewLogger()

	addr := ":53"

	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Temp local db connection
	db_conn, err := pgx.Connect(context.Background(), "postgresql://flip:postgres@localhost:5432/flipdns")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer conn.Close()
	defer db_conn.Close(ctx)

	repo := repository.New(db_conn)
	app := &app.App{
		Db:     repo,
		Logger: logger,
	}

	logger.Info("Started UDP server on addr " + addr)

	server := udpserver.NewUdpServer(app, conn)
	server.Run()

}

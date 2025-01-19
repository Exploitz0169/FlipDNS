package app

import (
	"log/slog"

	"github.com/exploitz0169/flipdns/internal/database"
)

type App struct {
	Db     *database.Database
	Logger *slog.Logger
}

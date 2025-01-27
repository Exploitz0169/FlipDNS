package app

import (
	"log/slog"

	"github.com/exploitz0169/flipdns/internal/repository"
)

type App struct {
	Db     *repository.Queries
	Logger *slog.Logger
}

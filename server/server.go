package server

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"log/slog"
	"os"

	"github.com/Route-E-106/Frogfoot/server/internal/database/models"
)

//go:embed internal/database/schema.sql
var ddl string

type Server struct {
	Logger  *slog.Logger
	Queries *models.Queries
	ctx     context.Context
}

func NewServer(db *sql.DB) *Server {

	ctx := context.Background()
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	return &Server{
		Logger:  slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Queries: models.New(db),
		ctx:     context.Background(),
	}

}

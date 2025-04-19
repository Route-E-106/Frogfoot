package server

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Route-E-106/Frogfoot/server/internal/database/models"
	"github.com/Route-E-106/Frogfoot/server/resources"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

//go:embed internal/database/schema.sql
var ddl string

type Server struct {
	Logger         *slog.Logger
	Queries        *models.Queries
	ctx            context.Context
	sessionManager *scs.SessionManager
	resources      resources.Resources
}

func NewServer(db *sql.DB) *Server {

	ctx := context.Background()
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 1 * time.Hour

	return &Server{
		Logger:         slog.New(slog.NewTextHandler(os.Stdout, nil)),
		Queries:        models.New(db),
		ctx:            context.Background(),
		sessionManager: sessionManager,
	}

}

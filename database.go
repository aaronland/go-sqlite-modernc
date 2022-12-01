package modernc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aaronland/go-sqlite/v2"
	_ "modernc.org/sqlite"
	"net/url"
	"sync"
	"log"
)

const SQLITE_SCHEME string = "modernc"
const SQLITE_DRIVER string = "sqlite"

type ModerncDatabase struct {
	sqlite.Database
	conn *sql.DB
	dsn  string
	mu   *sync.Mutex
	logger *log.Logger
}

func init() {
	ctx := context.Background()
	sqlite.RegisterDatabase(ctx, SQLITE_SCHEME, NewModerncDatabase)
}

func NewModerncDatabase(ctx context.Context, db_uri string) (sqlite.Database, error) {

	u, err := url.Parse(db_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	host := u.Host
	path := u.Path
	q := u.RawQuery

	var dsn string

	if host == "mem" {
		dsn = "file::memory:?mode=memory&cache=shared"
	} else {
		dsn = path
	}

	if q != "" {
		dsn = fmt.Sprintf("%s?%s", dsn, q)
	}

	conn, err := sql.Open(SQLITE_DRIVER, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	mu := new(sync.Mutex)

	logger := log.Default()
	
	db := ModerncDatabase{
		conn: conn,
		dsn:  dsn,
		mu:   mu,
		logger: logger,
	}

	return &db, nil
}

func (db *ModerncDatabase) Lock(ctx context.Context) error {
	db.mu.Lock()
	return nil
}

func (db *ModerncDatabase) Unlock(ctx context.Context) error {
	db.mu.Unlock()
	return nil
}

func (db *ModerncDatabase) Conn(ctx context.Context) (*sql.DB, error) {
	return db.conn, nil
}

func (db *ModerncDatabase) Close(ctx context.Context) error {
	return db.conn.Close()
}

func (db *ModerncDatabase) SetLogger(ctx context.Context, logger *log.Logger) error {
	db.logger = logger
	return nil
}

func (db *ModerncDatabase) DSN(ctx context.Context) string {
	return db.dsn
}

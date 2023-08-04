package sqlbase

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/utils"

	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Snippets contains embedded SQL files.
//
//go:embed snippets/*.sql
var SQL embed.FS

type SQLBase struct {
	driver string
	db     *sql.DB
	reset  bool
}

func NewSQLBase(cfg *models.Config) *SQLBase {
	var err error = nil
	sb := &SQLBase{
		driver: "pgx",
		reset:  cfg.Cleanup,
	}

	db, err := sql.Open("pgx", cfg.DatabaseDsn)

	if err != nil {
		panic(err)
	}
	sb.db = db

	logger.Log("Wait for DB reply... $")
	utils.Retryer(func() error {
		err = sb.db.Ping()
		return err
	})

	if err != nil {
		logger.Printf("FAILED\n")
		panic(err)
	}
	logger.Printf("OK\n")

	if sb.reset {
		err := sb.drop()
		if err != nil {
			logger.Error(err)
		}
	}

	err = sb.Init()
	if err != nil {
		logger.Error(err)
	}

	return sb
}

func (sb *SQLBase) drop() error {
	sqlQuery := "DROP SCHEMA IF EXISTS metmon CASCADE;"
	_, err := sb.db.Exec(sqlQuery)
	if err != nil {
		return fmt.Errorf("failed to drop schema: %w", err)
	}
	logger.Log("Schema dropped successfully")
	return nil
}

func (sb *SQLBase) Init() error {
	var err error = nil
	sqlBytes, err := SQL.ReadFile("snippets/pgdb_create_2.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}
	sqlQuery := string(sqlBytes)
	_, err = sb.db.Exec(sqlQuery)
	if err != nil {
		if strings.Contains(err.Error(), pgerrcode.DuplicateSchema) {
			logger.Log("Schema already exists")
			return nil
		}
	}
	logger.Log("Schema updated")
	return nil
}

func (sb *SQLBase) Close() {
	sb.db.Close()
}

func (sb *SQLBase) PingDB() error {
	err := sb.db.Ping()
	if err != nil {
		return fmt.Errorf("DB ping failed: %v", err.Error())
	}
	logger.Log("Successful DB ping")
	return nil
}

func loadSnippet(path string) (content string) {
	sqlBytes, err := SQL.ReadFile(path)
	if err != nil {
		logger.Error(fmt.Errorf("failed to read SQL file: %w", err))
		return ""
	}
	return string(sqlBytes)
}

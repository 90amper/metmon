package sqlbase

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"

	"github.com/90amper/metmon/internal/logger"
	"github.com/90amper/metmon/internal/models"
	"github.com/90amper/metmon/internal/utils"
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
	// ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	`localhost`, `video`, `XXXXXXXX`, `video`)

	// var db *sql.DB
	// utils.Retryer(func() error {
	// 	db, err = sql.Open("pgx", cfg.DatabaseDsn)
	// 	return err
	// })
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
	// err = sb.db.Ping()
	if err != nil {
		logger.Printf("FAILED\n")
		// return fmt.Errorf("DB ping failed: %v", err)
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
	// defer db.Close()
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
		// log.Error().Err(err).Msg("failed to read SQL file")
		return fmt.Errorf("failed to read SQL file: %w", err)
	}
	sqlQuery := string(sqlBytes)
	// sqlQuery := "SELECT 1+1"
	// fmt.Printf("%v; %v; %v", sb.db.Stats().MaxOpenConnections, sb.db.Stats().OpenConnections, sb.db.Stats().InUse)
	_, err = sb.db.Exec(sqlQuery)
	// if errors.Is(err, errors.New(pgerrcode.DuplicateSchema)) {
	if err != nil {
		if strings.Contains(err.Error(), pgerrcode.DuplicateSchema) {
			logger.Log("Schema already exists")
			return nil
		}
	}
	// logger.Log("Affcted rows: " + fmt.Sprint(aff))
	logger.Log("Schema updated")
	return nil
}

func (sb *SQLBase) Close() {
	sb.db.Close()
}

func (sb *SQLBase) PingDB() error {
	err := sb.db.Ping()
	// var greeting string
	// err := sb.db.QueryRow("select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return fmt.Errorf("DB ping failed: %v", err.Error())
		// os.Exit(1)
	}
	logger.Log("Successful DB ping")
	// fmt.Println(greeting)
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

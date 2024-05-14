package migrations

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type Migration struct {
	DB DB
}

type migration struct {
	FileName string
	Hash     string
}

func Run(ctx context.Context, db DB, path string) error {
	log.Info().Caller().Msg("Start run migrations")
	c := &Migration{DB: db}
	if err := c.init(ctx); err != nil {
		return err
	}

	if err := c.runMigration(ctx, path); err != nil {
		return err
	}
	log.Info().Caller().Msg("Finish run migrations")
	return nil
}

func (c *Migration) init(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS migrations (
        id SERIAL PRIMARY KEY,
        table_hash TEXT NOT NULL,
        table_name TEXT NOT NULL UNIQUE
    );
    `
	_, err := c.DB.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}
	return nil
}

func (c *Migration) getMigrations(ctx context.Context) (map[string]string, int, error) {
	rows, err := c.DB.Query(ctx, `SELECT table_name, table_hash FROM migrations ORDER BY table_name`)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying migrations: %v", err)
	}

	defer rows.Close()

	existingMigrations := make(map[string]string)
	var maxAppliedPrefix int
	for rows.Next() {
		var m migration
		if err := rows.Scan(&m.FileName, &m.Hash); err != nil {
			return nil, 0, fmt.Errorf("error scanning migration: %v", err)
		}
		existingMigrations[m.FileName] = m.Hash
		prefix, _ := strconv.Atoi(strings.Split(m.FileName, "_")[0])
		if prefix > maxAppliedPrefix {
			maxAppliedPrefix = prefix
		}
	}

	return existingMigrations, maxAppliedPrefix, nil
}

func (c *Migration) getMigrationsFiles(existingMigrations map[string]string, maxAppliedPrefix int, path string) ([]migration, error) {
	files, err := os.ReadDir(path)
	log.Info().Caller().Msg(fmt.Sprintf("migrations files: %s", files))
	if err != nil {
		return nil, fmt.Errorf("error reading migrations directory: %v", err)
	}

	var migrations []migration
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			prefix, err := strconv.Atoi(strings.Split(file.Name(), "_")[0])
			if err != nil || prefix < 0 {
				continue
			}
			if prefix <= maxAppliedPrefix && len(existingMigrations) > 0 {
				actualHash, exists := existingMigrations[file.Name()]
				if !exists {
					return nil, fmt.Errorf("invalid filename %s", file.Name())
				}
				content, err := os.ReadFile(filepath.Join(path, file.Name()))
				if err != nil {
					return nil, fmt.Errorf("error reading migration file %s: %v", file.Name(), err)
				}
				hash := md5.Sum(content)
				hashStr := hex.EncodeToString(hash[:])
				if actualHash != hashStr {
					return nil, fmt.Errorf("invalid hash migration file %s: %v", file.Name(), err)
				}
				continue
			}
			migrations = append(migrations, migration{FileName: file.Name()})
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		numI, _ := strconv.Atoi(strings.Split(migrations[i].FileName, "_")[0])
		numJ, _ := strconv.Atoi(strings.Split(migrations[j].FileName, "_")[0])
		return numI < numJ
	})

	return migrations, nil
}

func (c *Migration) runMigration(ctx context.Context, path string) error {
	existingMigrations, maxAppliedPrefix, err := c.getMigrations(ctx)
	if err != nil {
		return err
	}
	log.Info().Caller().Msg(fmt.Sprintf("existingMigrations: %s", existingMigrations))
	log.Info().Caller().Msg(fmt.Sprintf("maxAppliedPrefix: %d", maxAppliedPrefix))
	migrations, err := c.getMigrationsFiles(existingMigrations, maxAppliedPrefix, path)
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		if _, exists := existingMigrations[migration.FileName]; exists {
			continue
		}
		content, err := os.ReadFile(filepath.Join(path, migration.FileName))
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", migration.FileName, err)
		}
		hash := md5.Sum(content)
		hashStr := hex.EncodeToString(hash[:])
		if existingHash, exists := existingMigrations[migration.FileName]; exists && existingHash != hashStr {
			log.Info().Caller().Msg(fmt.Sprintf("migration %s exists", migration.FileName))
			continue
		}
		_, err = c.DB.Exec(context.Background(), string(content))
		if err != nil {
			return fmt.Errorf("error executing migration %s: %v", migration.FileName, err)
		}
		_, err = c.DB.Exec(context.Background(), `INSERT INTO migrations (table_hash, table_name) VALUES ($1, $2)`, hashStr, migration.FileName)
		if err != nil {
			return fmt.Errorf("error recording migration %s: %v", migration.FileName, err)
		}
		log.Info().Caller().Msg(fmt.Sprintf("migration %s executed", migration.FileName))
	}

	return nil
}

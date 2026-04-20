package rw

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/risingwavelabs/events-api/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestFilterInsertableColumns(t *testing.T) {
	cols := []Column{
		{Name: "id", Type: "integer"},
		{Name: "_row_id", Type: "integer"},
		{Name: "_rw_timestamp", Type: "timestamptz", IsHidden: true},
		{Name: "ingested_at", Type: "timestamptz", IsGenerated: true},
		{Name: "data", Type: "character varying"},
	}

	filtered := filterInsertableColumns(cols)
	require.Len(t, filtered, 2)
	require.Equal(t, []string{"id", "data"}, []string{filtered[0].Name, filtered[1].Name})
}

func TestParseEscapesPasswordUserinfo(t *testing.T) {
	password := "p@ss:/?#%[]word"
	dsn := parse(&config.Rw{
		Host:     "localhost",
		Port:     4566,
		User:     "root",
		Password: password,
		Db:       "dev",
		SSLMode:  "disable",
	})

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	require.Equal(t, "root", pgxCfg.ConnConfig.User)
	require.Equal(t, password, pgxCfg.ConnConfig.Password)
	require.Equal(t, "localhost", pgxCfg.ConnConfig.Host)
	require.Equal(t, uint16(4566), pgxCfg.ConnConfig.Port)
	require.Equal(t, "dev", pgxCfg.ConnConfig.Database)
}

func TestParseNormalizesInvalidDSNUserinfo(t *testing.T) {
	rawDSN := "postgres://root:p@ss:/?#%[]word@localhost:4566/dev?sslmode=disable"
	dsn := parse(&config.Rw{
		DSN: &rawDSN,
	})

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	require.Equal(t, "root", pgxCfg.ConnConfig.User)
	require.Equal(t, "p@ss:/?#%[]word", pgxCfg.ConnConfig.Password)
	require.Equal(t, "localhost", pgxCfg.ConnConfig.Host)
	require.Equal(t, uint16(4566), pgxCfg.ConnConfig.Port)
	require.Equal(t, "dev", pgxCfg.ConnConfig.Database)
}

func TestParseKeepsValidDSN(t *testing.T) {
	rawDSN := "postgres://root:p%40ss@localhost:4566/dev?sslmode=disable"

	require.Equal(t, rawDSN, parse(&config.Rw{
		DSN: &rawDSN,
	}))
}

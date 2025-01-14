package core_test

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/mysql"
	"github.com/orlangure/gnomock/preset/postgres"
)

type dbinfo struct {
	name    string
	driver  string
	connstr string
	disable bool
	preset  gnomock.Preset
}

var (
	dbParam string
	dbType  string
	db      *sql.DB
)

func init() {
	flag.StringVar(&dbParam, "db", "", "database type")
}

func TestMain(m *testing.M) {
	flag.Parse()

	dbinfoList := []dbinfo{
		{
			name:    "postgres",
			driver:  "postgres",
			connstr: "postgres://tester:tester@%s/db?sslmode=disable",
			preset: postgres.Preset(
				postgres.WithUser("tester", "tester"),
				postgres.WithDatabase("db"),
				postgres.WithQueriesFile("./postgres.sql"),
			)},
		{
			disable: true,
			name:    "mysql",
			driver:  "mysql",
			connstr: "user:user@tcp(%s)/db",
			preset: mysql.Preset(
				mysql.WithUser("user", "user"),
				mysql.WithDatabase("db"),
				mysql.WithQueriesFile("./mysql.sql"),
			)},
	}

	for _, v := range dbinfoList {
		disable := v.disable

		if dbParam != "" {
			if dbParam != v.name {
				continue
			} else {
				disable = false
			}
		}

		if disable {
			continue
		}

		con, err := gnomock.Start(v.preset)
		if err != nil {
			panic(err)
		}
		defer func() { _ = gnomock.Stop(con) }()

		db, err = sql.Open(v.driver, fmt.Sprintf(v.connstr, con.DefaultAddress()))
		if err != nil {
			panic(err)
		}
		db.SetMaxIdleConns(100)
		dbType = v.name

		if res := m.Run(); res != 0 {
			os.Exit(res)
		}
	}
}

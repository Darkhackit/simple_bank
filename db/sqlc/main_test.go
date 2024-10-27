package db

import (
	"context"
	"github.com/Darkhackit/simplebank/util"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	testDB, err = pgxpool.New(ctx, config.DbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer testDB.Close()
	testQueries = New(testDB)

	os.Exit(m.Run())

}

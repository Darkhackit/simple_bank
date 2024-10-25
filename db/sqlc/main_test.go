package db

import (
	"context"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *pgxpool.Pool

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	testDB, err = pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer testDB.Close()
	testQueries = New(testDB)

	os.Exit(m.Run())

}

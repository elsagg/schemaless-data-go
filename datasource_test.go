package schemaless

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-testfixtures/testfixtures"
)

type Info struct {
	FistName string `msgpack:"first_name"`
	LastName string `msgpack:"last_name"`
}

type Address struct {
	Street string `msgpack:"street"`
	Zip    int    `msgpack:"zip"`
}

var (
	db       *sql.DB
	fixtures *testfixtures.Context
	conn     = NewConnection(&ConnectionOptions{
		Host:       "127.0.0.1",
		Port:       "3306",
		Username:   "root",
		Password:   "",
		Protocol:   "tcp",
		Parameters: map[string]string{"parseTime": "true"},
	})
	ds = NewDataSource("users", conn)
)

func TestMain(m *testing.M) {
	var err error

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1)/users?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			added_id bigint(20) NOT NULL AUTO_INCREMENT,
			row_key varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
			column_key varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
			body blob,
			ref_key int(11) DEFAULT NULL,
			created_at datetime DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (added_id)
		  ) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	if err != nil {
		log.Fatal(err)
	}

	testfixtures.SkipDatabaseNameCheck(true)

	fixtures, err = testfixtures.NewFolder(db, &testfixtures.MySQL{}, "testdata/fixtures")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func TestGetCell(t *testing.T) {
	prepareTestDatabase()

	var info Info

	expected := "John"

	cell, err := ds.GetCell(
		"185133a4-71e4-4595-8f09-3ffe416064fc",
		"BASIC_INFO",
		1,
	)

	if err != nil {
		t.Error(err)
	}

	err = cell.UnmarshalBody(&info)

	if err != nil {
		t.Error(err)
	}

	if info.FistName != expected {
		t.Error(fmt.Sprintf("Expected FirstName to be %s, got %s", expected, info.FistName))
	}
}

func TestGetCellLatest(t *testing.T) {
	prepareTestDatabase()

	var address Address

	expected := 256255254

	cell, err := ds.GetCellLatest(
		"185133a4-71e4-4595-8f09-3ffe416064fc",
		"ADDRESS",
	)

	if err != nil {
		t.Error(err)
	}

	err = cell.UnmarshalBody(&address)

	if err != nil {
		t.Error(err)
	}

	if address.Zip != expected {
		t.Error(fmt.Sprintf("Expected ZIP code to be %d, got %d", expected, address.Zip))
	}
}

func TestGetAllLatest(t *testing.T) {
	prepareTestDatabase()

	expected := "185133a4-71e4-4595-8f09-3ffe416064fc"

	cells, err := ds.GetAllLatest("BASIC_INFO")

	if err != nil {
		t.Error(err)
	}

	if (*cells)[0].RowKey != expected {
		t.Error(fmt.Sprintf("Expected first address RowKey to be %s, got %s", expected, (*cells)[0].RowKey))
	}
}

func TestPutCell(t *testing.T) {
	prepareTestDatabase()

	var info Info

	expected := "Cesar"

	cell, err := ds.PutCell(
		"2860cd09-a7de-427e-8d5d-19353a84ddce",
		"BASIC_INFO",
		&Info{FistName: "Julio", LastName: "Cesar"},
	)

	if err != nil {
		t.Error(err)
	}

	err = cell.UnmarshalBody(&info)

	if err != nil {
		t.Error(err)
	}

	if info.LastName != expected {
		t.Error(fmt.Sprintf("Expected LastName to be %s, got %s", expected, info.LastName))
	}
}

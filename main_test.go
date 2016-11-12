package main

import (
	"errors"
	_ "github.com/DATA-DOG/go-sqlmock"
	"github.com/klichukb/analytics/client"
	"github.com/klichukb/analytics/server"
	"os"
	"testing"
	"time"
)

var (
	testDbName = "test_analytics"
	testDbUser = "analytics"
	testDbPwd  = "analytics"
)

// Tests interaction between client/server.
// Writes several messages from client to server and ensures events got written to DB.
func TestClientServer(t *testing.T) {
	// common mock flag parameters
	*address = ":8111"

	// server runner
	go runServer()

	// graceful wait for server
	time.Sleep(100 * time.Millisecond)
	// sets buffer of analytics object created by server
	analytics.MaxBufferSize = 4

	// client runner
	*workerCount = 2
	client.Iterations = 2
	runClient()

	row := server.DB.QueryRow("SELECT COUNT(*) FROM analytics_event")
	var rowCount int
	row.Scan(&rowCount)

	// NOTE: relies on fact that rows do get commited by this time.
	if rowCount != 4 {
		t.Errorf("Not all records were written to DB")
	}
}

// Sets up test database and tears down in the end.
// Clears data from events table.
func TestMain(m *testing.M) {
	// Setup
	// mock DB
	dbName := server.DbName
	dbUser := server.DbUser
	dbPwd := server.DbPwd

	server.DbName = testDbName
	server.DbUser = testDbUser
	server.DbPwd = testDbPwd

	// teardown test db settings
	defer func() {
		server.DbName = dbName
		server.DbUser = dbUser
		server.DbPwd = dbPwd
	}()

	// connect to test db
	server.InitDatabase()
	// close test db
	defer func() {
		server.DB.Close()
		server.DB = nil
	}()

	row := server.DB.QueryRow("SELECT DATABASE();")
	var ensureDbName string
	row.Scan(&ensureDbName)

	// must must must make sure
	if ensureDbName != testDbName {
		panic(errors.New("Failed to connect to test database"))
	}

	server.DB.Exec("DELETE FROM analytics_event")

	// run tests
	retCode := m.Run()
	os.Exit(retCode)
}

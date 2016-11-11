// Database handling for server side of analytics application.

package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/klichukb/analytics/shared"
	"os"
	"time"
)

// Credentials and target database for database connection.
// Make sure to mock these for test suite.
var (
	DbName   = os.Getenv("MYSQL_DB")
	DbUser   = os.Getenv("MYSQL_USER")
	DbPwd    = os.Getenv("MYSQL_PWD")
	DbDriver = "mysql"
)

// Connect to DB and return DB instance.
func GetDatabase() *sql.DB {
	conn := fmt.Sprintf("%s:%s@/%s", DbUser, DbPwd, DbName)
	db, err := sql.Open(DbDriver, conn)
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	return db
}

// Writes event object to database.
// Serializes params to JSON.
func SaveEvent(event *shared.Event) error {
	// Params are stores as JSON since amount and type of values in these
	// varies.
	jsonParams, err := json.Marshal(event.Params)

	if err != nil {
		return err
	}

	// ORM looks like overkill at the moment.
	_, dbErr := DB.Exec(`
        INSERT INTO analytics_event (event_type, ts, params)
        VALUES (?, ?, ?)
        `, event.EventType, time.Unix(int64(event.TS), 0), jsonParams,
	)

	if dbErr != nil {
		return dbErr
	}
	return nil
}

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

var (
	dbName = os.Getenv("MYSQL_DB")
	dbUser = os.Getenv("MYSQL_USER")
	dbPwd  = os.Getenv("MYSQL_PWD")
)

// Connect to DB and return DB instance.
func GetDatabase() *sql.DB {
	conn := fmt.Sprintf("%s:%s@/%s", dbUser, dbPwd, dbName)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		db.Close()
		panic(err.Error())
	}
	return db
}

// Writes event object to database.
// Serializes params to JSON.
func SaveEvent(event *shared.Event) error {
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

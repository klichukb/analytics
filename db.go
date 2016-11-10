package main

import "time"
import "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "encoding/json"
import "os"

var (
	dbName = os.Getenv("MYSQL_DB")
	dbUser = os.Getenv("MYSQL_USER")
	dbPwd  = os.Getenv("MYSQL_PWD")
)

// Connect to DB and return DB instance.
func getDatabase() *sql.DB {
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
func saveEvent(event *Event) error {
	jsonParams, err := json.Marshal(event.Params)

	if err != nil {
		return err
	}

	// ORM looks like overkill at the moment.
	_, dbErr := db.Exec(`
        INSERT INTO analytics_event (event_type, ts, params)
        VALUES (?, ?, ?)
        `, event.EventType, time.Unix(int64(event.TS), 0), jsonParams,
	)

	if dbErr != nil {
		return dbErr
	}
	return nil
}

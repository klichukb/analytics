package server

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/klichukb/analytics/shared"
	"testing"
	"time"
)

func setupDB() (sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	DB = db
	return mock, err
}

func getTestEvent() *shared.Event {
	ts := 1478818690
	params := map[string]interface{}{
		"key1": 123,
		"key2": "value2",
	}
	return &shared.Event{"session_start", ts, params}
}

func TestSaveEvent(t *testing.T) {
	mock, _ := setupDB()
	defer func() {
		DB.Close()
		DB = nil
	}()

	event := getTestEvent()
	tm := time.Unix(int64(event.TS), 0)

	mock.ExpectExec(
		`INSERT INTO analytics_event`,
	).WithArgs("session_start", tm, []byte(`{"key1":123,"key2":"value2"}`))

	SaveEvent(event)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations from `SaveEvent`: %s", err)
	}
}

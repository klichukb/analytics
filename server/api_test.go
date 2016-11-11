package server

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/klichukb/analytics/shared"
	"os"
	"testing"
)

var mock sqlmock.Sqlmock

func tryTrackEvent(event *shared.Event) error {
	analytics := &Analytics{}
	var reply int
	return analytics.TrackEvent(event, &reply)
}

func TestTrackEvent_valid(t *testing.T) {
	mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
	event := getTestEvent()
	err := tryTrackEvent(event)
	if err != nil {
		t.Error("Valid event should be successfully written")
	}
}

func TestTrackEvent_invalid(t *testing.T) {
	mock.ExpectExec("").WillReturnResult(sqlmock.NewResult(0, 0))
	event := getTestEvent()
	err := tryTrackEvent(event)
	if err != nil {
		t.Error("Valid event should be successfully written")
	}
}

func TestMain(m *testing.M) {
	dbMock, _ := setupDB()
	mock = dbMock

	defer func() {
		DB.Close()
		DB = nil
	}()

	retCode := m.Run()
	// call with result of m.Run()
	os.Exit(retCode)
}

package shared

type Event struct {
	EventType string
	TS        int
	Params    map[string]interface{}
}

var EventTypes = []string{
	"session_start",
	"session_end",
	"link_clicked",
}

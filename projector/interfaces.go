package projector

import "time"

type (
	DocumentMessage struct {
		Documents []Document
		Receipt   interface{}
	}

	Document interface {
		Lapse(now time.Time) (next Document)
		Handle(message interface{}) bool
		Path() string
	}
)

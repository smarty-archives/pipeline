package projector

import "time"

type (
	DocumentMessage struct {
		Documents []Document
		Receipt   interface{}
	}

	Document interface {
		Lapse(now time.Time) (next Document)
		Apply(message interface{}) bool
		Path() string
	}
)

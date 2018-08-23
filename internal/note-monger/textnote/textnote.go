package textnote

import "time"

type TextNote struct {
	Title        string
	Body         string
	DateCreated  time.Time
	DateModified time.Time
}

package entites

import "time"

type Appointment struct {
	ID            int64
	UserId        int64
	Services      []Service
	Date          time.Time
	Hour          string
	Minute        string
	TotalDuration int
}

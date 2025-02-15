package entites

import (
	"time"
)

type Appointment struct {
	ID            int
	UserId        int64
	Services      []Service
	Date          time.Time
	DateStr       string
	Hour          string
	Minute        string
	TotalDuration int
	TotalCost     int
}

package entites

import "time"

type Appointment struct {
	ID       int64
	UserId   int64
	DateTime time.Time
}

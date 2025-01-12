package rusdate

import (
	"fmt"
	"time"
)

var russianMonths = []string{
	"января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

func FormatDayMonth(t time.Time) string {
	day := t.Day()
	month := russianMonths[t.Month()-1]

	date := fmt.Sprintf("%d %s", day, month)
	return date
}

package rusdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var russianMonths = []string{
	"января", "февраля", "марта", "апреля", "мая", "июня",
	"июля", "августа", "сентября", "октября", "ноября", "декабря",
}

var russianMonthsMap = map[string]int{
	"января":   1,
	"февраля":  2,
	"марта":    3,
	"апреля":   4,
	"мая":      5,
	"июня":     6,
	"июля":     7,
	"августа":  8,
	"сентября": 9,
	"октября":  10,
	"ноября":   11,
	"декабря":  12,
}

func FormatDayMonth(t time.Time) string {
	day := t.Day()
	month := russianMonths[t.Month()-1]
	year := t.Year()

	date := fmt.Sprintf("%d %s %d", day, month, year)
	return date
}

func FormatBack(rusDate string) (time.Time, error) {
	params := strings.Split(rusDate, " ")

	day, err := strconv.Atoi(params[0])
	if err != nil {
		return time.Now(), err
	}

	month := russianMonthsMap[params[1]]

	year, err := strconv.Atoi(params[2])
	if err != nil {
		return time.Now(), err
	}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return date, nil
}

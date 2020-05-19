package qr

import "time"

type date time.Time

func (d date) String() string {
	formatted := time.Time(d).Format("20060102150405-07")
	// Trim the timezone to
	if formatted[15] == '0' {
		formatted = formatted[:15] + formatted[16:]
	}
	return formatted
}

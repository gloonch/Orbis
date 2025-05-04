package astro

import "time"

// ToJED converts a time.Time (UTC) into a Julian Ephemeris Date.
func ToJED(t time.Time) float64 {
	tu := t.UTC()

	year, month, day := tu.Date()
	hour, min, sec := tu.Clock()
	ns := tu.Nanosecond()

	Y := year
	M := int(month)
	D := float64(day)

	a := (14 - M) / 12
	y := Y + 4800 - a
	m := M + 12*a - 3

	JD0 := D +
		float64((153*m+2)/5) +
		365*float64(y) +
		float64(y/4) -
		float64(y/100) +
		float64(y/400) -
		32045

	frac := (float64(hour)-12.0)/24.0 +
		float64(min)/1440.0 +
		(float64(sec)+float64(ns)/1e9)/86400.0

	return JD0 + frac
}

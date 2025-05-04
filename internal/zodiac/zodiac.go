package zodiac

import (
	"math"
)

// EclipticLongitude computes the ecliptic longitude (in degrees, 0–360)
// of a body at rectangular coordinates (X,Y,Z) in the ecliptic plane.
// It ignores latitude (Z) for longitude calculation.
func EclipticLongitude(x, y, z float64) float64 {
	// atan2 returns radians between –π and +π
	rad := math.Atan2(y, x)
	deg := rad * 180 / math.Pi
	if deg < 0 {
		deg += 360
	}
	return deg
}

// Sign returns the zodiac sign name (Aries, Taurus, … Pisces)
// corresponding to the given ecliptic longitude in degrees.
func Sign(angle float64) string {
	signs := []string{
		"Aries", "Taurus", "Gemini", "Cancer",
		"Leo", "Virgo", "Libra", "Scorpio",
		"Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
	idx := int(angle/30) % 12
	return signs[idx]
}

// House returns the “equal house” number (1–12) for the given
// ecliptic longitude in degrees. House 1 = 0–30°, House 2 = 30–60°, etc.
func House(angle float64) int {
	h := int(angle/30) + 1
	if h > 12 {
		h = ((h-1)%12 + 1)
	}
	return h
}

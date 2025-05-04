package astro

import (
	"github.com/mshafiee/jpleph"
)

// GetBodyPosition loads the given ephemeris file (e.g. "de440.bin"),
// computes the position of `target` relative to `center` at JED `jed`,
// and returns only the 3-D position in AU.
func GetBodyPosition(
	ephFile string,
	jed float64,
	target jpleph.Planet,
	center jpleph.CenterBody,
) (jpleph.Position, error) {
	// open ephemeris; `false` means “don’t pre-load constants”
	eph, err := jpleph.NewEphemeris(ephFile, false)
	if err != nil {
		return jpleph.Position{}, err
	}
	defer eph.Close()

	// CalculatePV(et, target, center, calcVelocity)
	// calcVelocity=false → skip velocity calculation
	pos, _, err := eph.CalculatePV(jed, target, center, false)
	if err != nil {
		return jpleph.Position{}, err
	}
	return pos, nil
}

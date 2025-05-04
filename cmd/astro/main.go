package main

import (
	"fmt"
	"github.com/gloonch/orbis/internal/astro"
	"log"

	"github.com/mshafiee/jpleph"
)

func main() {
	ephFile := "ephemeris/linux_p1550p2650.440"
	jed := 2451545.0             // J2000.0
	target := jpleph.Mars        // any of jpleph.Mercury … jpleph.Sun, jpleph.Pluto…
	center := jpleph.CenterEarth // or CenterSun, CenterMoon… for any reference point

	pos, err := astro.GetBodyPosition(ephFile, jed, target, center)
	if err != nil {
		log.Fatalf("error getting position: %v", err)
	}

	fmt.Printf(
		"At JED %.1f, %v relative to %v is at X=%.7f, Y=%.7f, Z=%.7f AU\n",
		jed, target, center, pos.X, pos.Y, pos.Z,
	)
}

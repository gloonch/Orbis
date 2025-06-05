// internal/astro/stream.go
package astro

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gloonch/orbis/internal/kafka"
	"github.com/gloonch/orbis/internal/zodiac" // Add this import
	"github.com/mshafiee/jpleph"
)

type PositionMessage struct {
	Body      string  `json:"body"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Z         float64 `json:"z"`
	Time      float64 `json:"time"`
	House     int     `json:"house"`
	Degree    float64 `json:"degree"`
	Longitude float64 `json:"longitude"` // New field
}

var planetNames = map[jpleph.Planet]string{
	jpleph.Sun:     "Sun",
	jpleph.Moon:    "Moon",
	jpleph.Mercury: "Mercury",
	jpleph.Venus:   "Venus",
	jpleph.Mars:    "Mars",
	jpleph.Jupiter: "Jupiter",
	jpleph.Saturn:  "Saturn",
	jpleph.Uranus:  "Uranus",
	jpleph.Neptune: "Neptune",
	jpleph.Pluto:   "Pluto",
}

func StartPlanetStream(
	ctx context.Context,
	ephFile string,
	body jpleph.Planet,
	prod *kafka.Producer,
	topic string,
	intervalSeconds int,
) {
	eph, err := jpleph.NewEphemeris(ephFile, false)
	if err != nil {
		log.Fatalf("failed to load ephemeris %s: %v", ephFile, err)
	}
	defer eph.Close()

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:

			jed := ToJED(t)

			pos, err := GetBodyPosition(ephFile, jed, body, jpleph.CenterEarth)
			if err != nil {
				log.Fatalf("error getting position: %v", err)
			}

			// Calculate ecliptic longitude, house, and degree
			longitude := zodiac.EclipticLongitude(pos.X, pos.Y, pos.Z)
			house := zodiac.House(longitude)
			degree := longitude / 30

			// 3️⃣ ساخت پیام و مارشال شدن به JSON
			msg := PositionMessage{
				Body:      planetNames[body],
				X:         pos.X,
				Y:         pos.Y,
				Z:         pos.Z,
				Time:      jed,
				House:     house,
				Degree:    degree,
				Longitude: longitude, // Assign longitude here
			}
			b, err := json.Marshal(msg)
			if err != nil {
				log.Printf("json marshal error: %v", err)
				continue
			}

			if err := prod.Publish(topic, planetNames[body], b); err != nil {
				log.Printf("publish error: %v", err)
			}
		}
	}
}

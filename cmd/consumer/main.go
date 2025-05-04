// cmd/consumer/main.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gloonch/orbis/config"
	"github.com/gloonch/orbis/internal/kafka"
	"github.com/gloonch/orbis/internal/zodiac"
)

// PositionMsg matches the JSON payload produced by the producer:
// { "X": ..., "Y": ..., "Z": ... }
type PositionMsg struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
	Z float64 `json:"Z"`
}

func main() {
	//  Load config
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	//  Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGINT/SIGTERM for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	// Create Kafka consumer
	cons, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.TopicPositions, // e.g. "positions"
		cfg.Kafka.GroupID,        // unique consumer group
	)
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}
	defer cons.Close()

	log.Println("consumer started, awaiting messages…")

	//  Start consuming
	err = cons.Consume(ctx, func(key string, val []byte) {
		// key is the planet name, e.g. "Mars"
		var msg PositionMsg
		if err := json.Unmarshal(val, &msg); err != nil {
			log.Printf("invalid message for %s: %v", key, err)
			return
		}

		//  Compute ecliptic longitude (degrees) from (X,Y,Z)
		angle := zodiac.EclipticLongitude(msg.X, msg.Y, msg.Z)

		// Determine zodiac sign and house
		sign := zodiac.Sign(angle)
		house := zodiac.House(angle)

		// Output result
		log.Printf(
			"%s → %.3f° → %s (House %d)",
			key, angle, sign, house,
		)
	})
	if err != nil && err != context.Canceled {
		log.Fatalf("consume error: %v", err)
	}

	log.Println("consumer shutting down")
}

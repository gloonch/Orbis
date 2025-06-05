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

// PositionMsg matches the JSON payload produced by the producer
type PositionMsg struct {
	Body      string  `json:"body"`
	X         float64 `json:"x"` // Corresponds to producer's json:"x"
	Y         float64 `json:"y"` // Corresponds to producer's json:"y"
	Z         float64 `json:"z"` // Corresponds to producer's json:"z"
	Time      float64 `json:"time"`
	House     int     `json:"house"`
	Degree    float64 `json:"degree"`
	Longitude float64 `json:"longitude"`
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
		// key is the planet name, e.g. "Mars" - msg.Body will be used instead
		var msg PositionMsg
		if err := json.Unmarshal(val, &msg); err != nil {
			log.Printf("invalid message for key %s: %v", key, err) // Keep key for error logging context
			return
		}

		// Determine zodiac sign using the Longitude from the message
		sign := zodiac.Sign(msg.Longitude)

		// Output result using fields directly from the message
		log.Printf(
			"%s → %.3f° %s (House %d)", // Format degree, then sign
			msg.Body, msg.Degree, sign, msg.House,
		)
	})
	if err != nil && err != context.Canceled {
		log.Fatalf("consume error: %v", err)
	}

	log.Println("consumer shutting down")
}

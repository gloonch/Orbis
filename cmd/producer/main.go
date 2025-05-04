package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gloonch/orbis/config"
	"github.com/gloonch/orbis/internal/astro"
	"github.com/gloonch/orbis/internal/kafka"
	"github.com/mshafiee/jpleph"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	topic := cfg.Kafka.TopicPositions

	// Create AdminClient
	admin, err := kafka.NewAdminClient(cfg.Kafka.Brokers[0])
	if err != nil {
		log.Fatalf("failed to create Kafka admin client: %v", err)
	}
	defer admin.Close()

	// Ensure the topic exists
	err = admin.CreateTopic(topic, 1, 1) // 1 partition, 1 replication factor
	if err != nil {
		log.Fatalf("failed to create Kafka topic: %v", err)
	}

	// Create Producer
	prod, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer prod.Close()

	planets := []jpleph.Planet{
		jpleph.Sun,
		jpleph.Moon,
		jpleph.Mercury,
		jpleph.Venus,
		jpleph.Mars,
		jpleph.Jupiter,
		jpleph.Saturn,
		jpleph.Uranus,
		jpleph.Neptune,
		jpleph.Pluto,
	}
	for _, planet := range planets {
		go astro.StartPlanetStream(
			ctx,
			cfg.Ephemeris.FilePath,
			planet,
			prod,
			topic,
			cfg.Stream.IntervalSeconds,
		)
	}

	<-ctx.Done()
	log.Println("shutting down producer gracefully…")
}

// internal/kafka/producer.go
package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// Producer is a simple Kafka producer wrapper.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Producer that writes to the given brokers.
func NewProducer(brokers []string) (*Producer, error) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		// we’ll write to whichever topic is specified at Publish time
		// so leave Topic empty here
		Balancer: &kafka.LeastBytes{},
		// you can tune these as needed
		BatchTimeout: 10 * time.Millisecond,
	})
	return &Producer{writer: w}, nil
}

// Publish sends a single message with the given key and value
// to the specified topic.
func (p *Producer) Publish(topic, key string, value []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	}
	// use a background context; you could also accept a ctx parameter
	return p.writer.WriteMessages(context.Background(), msg)
}

// Close flushes and closes the underlying writer.
func (p *Producer) Close() error {
	return p.writer.Close()
}

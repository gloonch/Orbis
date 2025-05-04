// internal/kafka/consumer.go
package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// Consumer is a simple Kafka consumer wrapper.
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a new Consumer that reads from the given brokers and topic,
// using the given groupID for consumer-group coordination.
func NewConsumer(brokers []string, topic, groupID string) (*Consumer, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{reader: r}, nil
}

// Consume starts pulling messages until ctx is done.
// For each message, it calls handler(key, value).
// It commits offsets after handler returns.
func (c *Consumer) Consume(ctx context.Context, handler func(key string, value []byte)) error {
	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			// context cancellation is expected on shutdown
			if err == context.Canceled {
				return nil
			}
			return err
		}

		fmt.Printf(
			"Consumed message topic=%s partition=%d offset=%d key=%q value=%q\n",
			m.Topic, m.Partition, m.Offset,
			string(m.Key), string(m.Value),
		)

		// invoke user handler
		handler(string(m.Key), m.Value)

		// commit offset so we won't re-process on restart
		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("failed to commit message: %v", err)
		}
	}
}

// Close closes the underlying reader.
func (c *Consumer) Close() error {
	return c.reader.Close()
}

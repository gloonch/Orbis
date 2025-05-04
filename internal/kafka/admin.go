// internal/kafka/admin.go
package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

// AdminClient is a simple Kafka admin client wrapper.
type AdminClient struct {
	conn *kafka.Conn
}

// NewAdminClient creates a new AdminClient connected to the given broker.
func NewAdminClient(broker string) (*AdminClient, error) {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return nil, err
	}
	return &AdminClient{conn: conn}, nil
}

// CreateTopic creates a new topic with the given name, partitions, and replication factor.
func (a *AdminClient) CreateTopic(topic string, partitions, replicationFactor int) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return a.conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})
}

// Close closes the admin client connection.
func (a *AdminClient) Close() error {
	return a.conn.Close()
}

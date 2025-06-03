package astro

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"math"

	"github.com/gloonch/orbis/internal/kafka"
	"github.com/mshafiee/jpleph"
	"github.com/stretchr/testify/assert"
)

// Mock GetBodyPosition for testing purposes
// It allows us to return a fixed position
func mockGetBodyPosition(ephFile string, jed float64, target jpleph.Planet, center jpleph.CenterBody) (jpleph.Position, error) {
	// For this test, we'll return a fixed position that should result in
	// longitude 0, house 1, degree 0.
	// X=1, Y=0, Z=0 corresponds to 0 degrees ecliptic longitude.
	return jpleph.Position{X: 1, Y: 0, Z: 0}, nil
}

// MockKafkaProducer is a mock implementation of the Kafka producer
type MockKafkaProducer struct {
	PublishedMessages []struct {
		Topic string
		Key   string
		Value []byte
	}
}

func (m *MockKafkaProducer) Publish(topic string, key string, value []byte) error {
	m.PublishedMessages = append(m.PublishedMessages, struct {
		Topic string
		Key   string
		Value []byte
	}{topic, key, value})
	return nil
}

func (m *MockKafkaProducer) Close() {}


func TestStartPlanetStream_CalculatesHouseAndDegree(t *testing.T) {
	// Replace the original GetBodyPosition with our mock for the duration of this test
	originalGetBodyPosition := GetBodyPosition
	GetBodyPosition = mockGetBodyPosition
	defer func() { GetBodyPosition = originalGetBodyPosition }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockProducer := &MockKafkaProducer{}
	ephFile := "test_ephemeris_file.eph" // Dummy file, as GetBodyPosition is mocked
	body := jpleph.Sun
	topic := "test_topic"
	intervalSeconds := 1 // Publish quickly for the test

	go StartPlanetStream(ctx, ephFile, body, mockProducer, topic, intervalSeconds)

	// Allow some time for at least one message to be produced
	time.Sleep(2 * time.Second)
	cancel() // Stop the streamer

	// Check if any message was published
	assert.NotEmpty(t, mockProducer.PublishedMessages, "No messages were published")

	// Unmarshal the first published message
	var msg PositionMessage
	err := json.Unmarshal(mockProducer.PublishedMessages[0].Value, &msg)
	assert.NoError(t, err, "Failed to unmarshal message")

	// Assertions for House and Degree
	// For X=1, Y=0, Z=0:
	// longitude = atan2(0,1) * 180/pi = 0
	// house = int(0/30) + 1 = 1
	// degree = 0 % 30 = 0
	expectedHouse := 1
	expectedDegree := 0.0

	assert.Equal(t, planetNames[body], msg.Body, "Message body mismatch")
	assert.Equal(t, expectedHouse, msg.House, "House calculation is incorrect")
	assert.InDelta(t, expectedDegree, msg.Degree, 0.0001, "Degree calculation is incorrect")

	// Verify X, Y, Z are as expected from the mock
	assert.Equal(t, 1.0, msg.X, "X coordinate mismatch")
	assert.Equal(t, 0.0, msg.Y, "Y coordinate mismatch")
	assert.Equal(t, 0.0, msg.Z, "Z coordinate mismatch")
}

package events

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EventHandler is called for each received event
type EventHandler func(ctx context.Context, event interface{}) error

// StreamClient defines the interface for event stream clients
type StreamClient interface {
	Recv() (interface{}, error)
}

// ConsumerConfig holds consumer configuration
type ConsumerConfig struct {
	Name           string
	DialFunc       func() (*grpc.ClientConn, error)
	StreamFunc     func(*grpc.ClientConn) (StreamClient, error)
	Handler        EventHandler
	Logger         *slog.Logger
	ReconnectDelay time.Duration
	MaxReconnect   time.Duration
}

// Consumer manages event stream consumption with automatic reconnection
type Consumer struct {
	name           string
	dialFunc       func() (*grpc.ClientConn, error)
	streamFunc     func(*grpc.ClientConn) (StreamClient, error)
	handler        EventHandler
	logger         *slog.Logger
	reconnectDelay time.Duration
	maxReconnect   time.Duration
	lastSeenSeq    map[string]uint64
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewConsumer creates a new event consumer with reconnection
func NewConsumer(cfg ConsumerConfig) *Consumer {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if cfg.ReconnectDelay == 0 {
		cfg.ReconnectDelay = 1 * time.Second
	}
	if cfg.MaxReconnect == 0 {
		cfg.MaxReconnect = 60 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		name:           cfg.Name,
		dialFunc:       cfg.DialFunc,
		streamFunc:     cfg.StreamFunc,
		handler:        cfg.Handler,
		logger:         cfg.Logger,
		reconnectDelay: cfg.ReconnectDelay,
		maxReconnect:   cfg.MaxReconnect,
		lastSeenSeq:    make(map[string]uint64),
		ctx:            ctx,
		cancel:         cancel,
	}
}

// Start begins consuming events with automatic reconnection
func (c *Consumer) Start(_ context.Context) error {
	c.logger.Info("starting event consumer", "name", c.name)
	go c.consumeLoop()
	return nil
}

// Stop gracefully stops the consumer
func (c *Consumer) Stop(_ context.Context) error {
	c.logger.Info("stopping event consumer", "name", c.name)
	c.cancel()
	return nil
}

// Name returns the consumer name
func (c *Consumer) Name() string {
	return c.name
}

// consumeLoop handles the reconnection logic
func (c *Consumer) consumeLoop() {
	backoff := c.reconnectDelay

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Info("consumer stopped", "name", c.name)
			return
		default:
		}

		// Attempt connection
		conn, err := c.dialFunc()
		if err != nil {
			c.logger.Error("failed to dial", "name", c.name, "error", err)
			c.sleep(backoff)
			backoff = c.increaseBackoff(backoff)
			continue
		}

		// Create stream
		stream, err := c.streamFunc(conn)
		if err != nil {
			c.logger.Error("failed to create stream", "name", c.name, "error", err)
			conn.Close()
			c.sleep(backoff)
			backoff = c.increaseBackoff(backoff)
			continue
		}

		// Consume events
		c.logger.Info("event stream connected", "name", c.name)
		backoff = c.reconnectDelay // Reset backoff on successful connection

		if err := c.consumeStream(stream); err != nil {
			c.logger.Error("stream error", "name", c.name, "error", err)
		}

		conn.Close()
		c.logger.Warn("stream disconnected, reconnecting", "name", c.name)
		c.sleep(backoff)
		backoff = c.increaseBackoff(backoff)
	}
}

// consumeStream processes events from the stream
func (c *Consumer) consumeStream(stream StreamClient) error {
	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		default:
		}

		event, err := stream.Recv()
		if err == io.EOF {
			c.logger.Info("stream closed by server", "name", c.name)
			return nil
		}
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.Canceled {
					return nil
				}
			}
			return fmt.Errorf("recv error: %w", err)
		}

		// Handle the event
		if err := c.handler(c.ctx, event); err != nil {
			c.logger.Error("handler error", "name", c.name, "error", err)
			// Continue processing despite handler errors
		}
	}
}

// sleep sleeps for the specified duration or until context is cancelled
func (c *Consumer) sleep(d time.Duration) {
	select {
	case <-c.ctx.Done():
	case <-time.After(d):
	}
}

// increaseBackoff increases the backoff duration with exponential backoff
func (c *Consumer) increaseBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > c.maxReconnect {
		return c.maxReconnect
	}
	return next
}

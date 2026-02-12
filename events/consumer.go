package events
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

// Consumer manages event stream consumption with automatic reconnection
type Consumer struct {
	name            string
	dialFunc        func() (*grpc.ClientConn, error)

































































































































































}	return next	}		return c.maxReconnect	if next > c.maxReconnect {	next := current * 2func (c *Consumer) increaseBackoff(current time.Duration) time.Duration {// increaseBackoff increases the backoff duration with exponential backoff}	}	case <-time.After(d):	case <-c.ctx.Done():	select {func (c *Consumer) sleep(d time.Duration) {// sleep sleeps for the specified duration or until context is cancelled}	}		}			// Continue processing despite handler errors			c.logger.Error("handler error", "name", c.name, "error", err)		if err := c.handler(c.ctx, event); err != nil {		// Handle the event		}			return fmt.Errorf("recv error: %w", err)			}				}					return nil				if st.Code() == codes.Canceled {			if st, ok := status.FromError(err); ok {		if err != nil {		}			return nil			c.logger.Info("stream closed by server", "name", c.name)		if err == io.EOF {		event, err := stream.Recv()		}		default:			return c.ctx.Err()		case <-c.ctx.Done():		select {	for {func (c *Consumer) consumeStream(stream StreamClient) error {// consumeStream processes events from the stream}	}		backoff = c.increaseBackoff(backoff)		c.sleep(backoff)		c.logger.Warn("stream disconnected, reconnecting", "name", c.name)		conn.Close()		}			c.logger.Error("stream error", "name", c.name, "error", err)		if err := c.consumeStream(stream); err != nil {		// Consume events		backoff = c.reconnectDelay // Reset backoff on successful connection		c.logger.Info("event stream connected", "name", c.name)		}			continue			backoff = c.increaseBackoff(backoff)			c.sleep(backoff)			conn.Close()			c.logger.Error("failed to create stream", "name", c.name, "error", err)		if err != nil {		stream, err := c.streamFunc(conn)		// Create stream		}			continue			backoff = c.increaseBackoff(backoff)			c.sleep(backoff)			c.logger.Error("failed to dial", "name", c.name, "error", err)		if err != nil {		conn, err := c.dialFunc()		// Attempt connection		}		default:			return			c.logger.Info("consumer stopped", "name", c.name)		case <-c.ctx.Done():		select {	for {	backoff := c.reconnectDelayfunc (c *Consumer) consumeLoop() {// consumeLoop handles the reconnection logic}	return c.namefunc (c *Consumer) Name() string {// Name returns the consumer name}	return nil	c.cancel()	c.logger.Info("stopping event consumer", "name", c.name)func (c *Consumer) Stop() error {// Stop gracefully stops the consumer}	return nil	go c.consumeLoop()	c.logger.Info("starting event consumer", "name", c.name)func (c *Consumer) Start() error {// Start begins consuming events with automatic reconnection}	}		cancel:         cancel,		ctx:            ctx,		lastSeenSeq:    make(map[string]uint64),		maxReconnect:   cfg.MaxReconnect,		reconnectDelay: cfg.ReconnectDelay,		logger:         cfg.Logger,		handler:        cfg.Handler,		streamFunc:     cfg.StreamFunc,		dialFunc:       cfg.DialFunc,		name:           cfg.Name,	return &Consumer{	ctx, cancel := context.WithCancel(context.Background())	}		cfg.MaxReconnect = 60 * time.Second	if cfg.MaxReconnect == 0 {	}		cfg.ReconnectDelay = 1 * time.Second	if cfg.ReconnectDelay == 0 {	}		cfg.Logger = slog.Default()	if cfg.Logger == nil {func NewConsumer(cfg ConsumerConfig) *Consumer {// NewConsumer creates a new event consumer with reconnection}	MaxReconnect   time.Duration	ReconnectDelay time.Duration	Logger         *slog.Logger	Handler        EventHandler	StreamFunc     func(*grpc.ClientConn) (StreamClient, error)	DialFunc       func() (*grpc.ClientConn, error)	Name           stringtype ConsumerConfig struct {// ConsumerConfig holds consumer configuration}	cancel          context.CancelFunc	ctx             context.Context	lastSeenSeq     map[string]uint64	maxReconnect    time.Duration	reconnectDelay  time.Duration	logger          *slog.Logger	handler         EventHandler	streamFunc      func(*grpc.ClientConn) (StreamClient, error)
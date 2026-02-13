package telemetry

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Event represents a telemetry or audit event
type Event struct {
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"`
	Component string                 `json:"component"`
	Action    string                 `json:"action"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// Flusher defines the interface for flushing telemetry
type Flusher interface {
	Flush(ctx context.Context, events []Event) error
}

// Buffer provides thread-safe buffering of telemetry events
type Buffer struct {
	mu      sync.Mutex
	events  []Event
	maxSize int
}

// NewBuffer creates a new telemetry buffer
func NewBuffer(maxSize int) *Buffer {
	return &Buffer{
		events:  make([]Event, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add appends an event to the buffer, dropping the oldest if full
func (b *Buffer) Add(event Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.events) >= b.maxSize {
		// Drop oldest event
		b.events = b.events[1:]
	}
	b.events = append(b.events, event)
}

// Drain removes and returns all buffered events
func (b *Buffer) Drain() []Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	events := b.events
	b.events = make([]Event, 0, b.maxSize)
	return events
}

// Len returns the number of buffered events
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.events)
}

// Collector manages periodic telemetry flushing
type Collector struct {
	name          string
	buffer        *Buffer
	flusher       Flusher
	flushInterval time.Duration
	logger        *slog.Logger
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewCollector creates a new telemetry collector
func NewCollector(name string, buffer *Buffer, flusher Flusher, flushInterval time.Duration, logger *slog.Logger) *Collector {
	if logger == nil {
		logger = slog.Default()
	}
	if flushInterval == 0 {
		flushInterval = 30 * time.Second
	}

	return &Collector{
		name:          name,
		buffer:        buffer,
		flusher:       flusher,
		flushInterval: flushInterval,
		logger:        logger,
	}
}

// Start begins the periodic flush loop
func (c *Collector) Start(_ context.Context) error {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.logger.Info("starting telemetry collector", "name", c.name, "interval", c.flushInterval)
	go c.flushLoop()
	return nil
}

// Stop gracefully stops the collector and performs a final flush
func (c *Collector) Stop(ctx context.Context) error {
	c.logger.Info("stopping telemetry collector", "name", c.name)
	c.cancel()

	// Final flush
	if err := c.FlushNow(ctx); err != nil {
		c.logger.Error("final flush failed", "name", c.name, "error", err)
		return err
	}
	return nil
}

// Name returns the collector name
func (c *Collector) Name() string {
	return c.name
}

// Record adds a telemetry event to the buffer
func (c *Collector) Record(eventType, action string, details map[string]interface{}) {
	c.buffer.Add(Event{
		Timestamp: time.Now(),
		Type:      eventType,
		Component: c.name,
		Action:    action,
		Details:   details,
	})
}

// FlushNow immediately flushes all buffered events
func (c *Collector) FlushNow(ctx context.Context) error {
	events := c.buffer.Drain()
	if len(events) == 0 {
		return nil
	}

	c.logger.Debug("flushing telemetry events", "name", c.name, "count", len(events))

	if err := c.flusher.Flush(ctx, events); err != nil {
		c.logger.Error("flush failed", "name", c.name, "error", err)
		return err
	}

	c.logger.Debug("telemetry events flushed", "name", c.name, "count", len(events))
	return nil
}

// flushLoop periodically flushes buffered events
func (c *Collector) flushLoop() {
	ticker := time.NewTicker(c.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.FlushNow(c.ctx); err != nil {
				c.logger.Error("periodic flush failed", "name", c.name, "error", err)
			}
		}
	}
}

package telemetry
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
	mu     sync.Mutex
	events []Event
	maxSize int
}

// NewBuffer creates a new telemetry buffer
func NewBuffer(maxSize int) *Buffer {
	return &Buffer{


























































































































































}	return nil	}		return err		}			c.buffer.Add(event)		for _, event := range events {		// Re-add events to buffer on failure	if err := c.flusher.Flush(ctx, events); err != nil {		c.logger.Debug("flushing telemetry", "name", c.name, "count", len(events))	}		return nil	if len(events) == 0 {	events := c.buffer.Drain()func (c *Collector) flush(ctx context.Context) error {// flush drains the buffer and sends events to the flusher}	}		}			}				c.logger.Error("flush failed", "name", c.name, "error", err)			if err := c.flush(c.ctx); err != nil {		case <-ticker.C:			return		case <-c.ctx.Done():		select {	for {	defer ticker.Stop()	ticker := time.NewTicker(c.flushInterval)	defer c.wg.Done()func (c *Collector) flushLoop() {// flushLoop periodically flushes events}	return c.flush(ctx)func (c *Collector) FlushNow(ctx context.Context) error {// FlushNow immediately flushes buffered events}	c.buffer.Add(event)	}		Details:   details,		Action:    action,		Component: c.name,		Type:      eventType,		Timestamp: time.Now(),	event := Event{func (c *Collector) Record(eventType, action string, details map[string]interface{}) {// Record adds an event to the buffer}	return c.namefunc (c *Collector) Name() string {// Name returns the collector name}	return nil	}		return err		c.logger.Error("final flush failed", "name", c.name, "error", err)	if err := c.flush(ctx); err != nil {	// Final flush	c.wg.Wait()	c.cancel()	c.logger.Info("stopping telemetry collector", "name", c.name)func (c *Collector) Stop(ctx context.Context) error {// Stop stops the collector and flushes remaining events}	return nil		go c.flushLoop()	c.wg.Add(1)		c.logger.Info("starting telemetry collector", "name", c.name, "interval", c.flushInterval)func (c *Collector) Start(ctx context.Context) error {// Start begins periodic flushing}	}		cancel:        cancel,		ctx:           ctx,		logger:        logger,		flushInterval: flushInterval,		flusher:       flusher,		buffer:        buffer,		name:          name,	return &Collector{	ctx, cancel := context.WithCancel(context.Background())	}		logger = slog.Default()	if logger == nil {func NewCollector(name string, buffer *Buffer, flusher Flusher, flushInterval time.Duration, logger *slog.Logger) *Collector {// NewCollector creates a new telemetry collector}	wg             sync.WaitGroup	cancel         context.CancelFunc	ctx            context.Context	logger         *slog.Logger	flushInterval  time.Duration	flusher        Flusher	buffer         *Buffer	name           stringtype Collector struct {// Collector manages periodic telemetry collection and flushing}	return len(b.events)	defer b.mu.Unlock()	b.mu.Lock()func (b *Buffer) Len() int {// Len returns the number of buffered events}	return drained	b.events = b.events[:0] // Reset slice	copy(drained, b.events)	drained := make([]Event, len(b.events))	defer b.mu.Unlock()	b.mu.Lock()func (b *Buffer) Drain() []Event {// Drain removes and returns all buffered events}	}		b.events = b.events[excess:]		excess := len(b.events) - b.maxSize		// Drop oldest events	if len(b.events) > b.maxSize {	// Prevent unbounded growth		b.events = append(b.events, event)	defer b.mu.Unlock()	b.mu.Lock()func (b *Buffer) Add(event Event) {// Add adds an event to the buffer}	}		maxSize: maxSize,		events:  make([]Event, 0, maxSize),
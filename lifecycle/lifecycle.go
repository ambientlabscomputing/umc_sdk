package lifecycle

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ambientlabscomputing/umc_sdk/types"
)

// Component represents a managed component with lifecycle methods
type Component interface {
	// Name returns the component name
	Name() string
	// Start initializes and starts the component
	Start(ctx context.Context) error
	// Stop gracefully shuts down the component
	Stop(ctx context.Context) error
}

// Launcher manages ordered initialization and teardown of components
type Launcher struct {
	name       string
	logger     *slog.Logger
	components []Component
	state      types.ComponentState
	mu         sync.RWMutex
}

// NewLauncher creates a new component launcher
func NewLauncher(logger *slog.Logger) *Launcher {
	if logger == nil {
		logger = slog.Default()
	}

	return &Launcher{
		name:       "launcher",
		logger:     logger,
		components: make([]Component, 0),
		state:      types.ComponentStateStopped,
	}
}

// Add registers a component for lifecycle management
func (l *Launcher) Add(component Component) {
	l.components = append(l.components, component)
}

// Start initializes all components in registration order
func (l *Launcher) Start(ctx context.Context) error {
	l.mu.Lock()
	l.state = types.ComponentStateStarting
	l.mu.Unlock()

	l.logger.Info("starting components", "count", len(l.components))

	for _, comp := range l.components {
		l.logger.Debug("starting component", "name", comp.Name())

		if err := comp.Start(ctx); err != nil {
			l.logger.Error("component start failed", "name", comp.Name(), "error", err)

			// Rollback: stop already started components in reverse order
			l.logger.Warn("rolling back started components")
			l.stopComponents(ctx, l.components[:indexOf(l.components, comp)])

			l.mu.Lock()
			l.state = types.ComponentStateFailed
			l.mu.Unlock()

			return fmt.Errorf("failed to start component %s: %w", comp.Name(), err)
		}

		l.logger.Info("component started", "name", comp.Name())
	}

	l.mu.Lock()
	l.state = types.ComponentStateRunning
	l.mu.Unlock()

	l.logger.Info("all components started successfully")
	return nil
}

// Stop shuts down all components in reverse order
func (l *Launcher) Stop(ctx context.Context) error {
	l.mu.Lock()
	l.state = types.ComponentStateStopping
	l.mu.Lock()

	l.logger.Info("stopping components", "count", len(l.components))

	errs := l.stopComponents(ctx, l.components)

	l.mu.Lock()
	l.state = types.ComponentStateStopped
	l.mu.Unlock()

	if len(errs) > 0 {
		l.logger.Error("errors during shutdown", "count", len(errs))
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	l.logger.Info("all components stopped successfully")
	return nil
}

// State returns the current launcher state
func (l *Launcher) State() types.ComponentState {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.state
}

// stopComponents stops components in reverse order and collects errors
func (l *Launcher) stopComponents(ctx context.Context, components []Component) []error {
	var errs []error

	// Stop in reverse order
	for i := len(components) - 1; i >= 0; i-- {
		comp := components[i]
		l.logger.Debug("stopping component", "name", comp.Name())

		if err := comp.Stop(ctx); err != nil {
			l.logger.Error("component stop failed", "name", comp.Name(), "error", err)
			errs = append(errs, fmt.Errorf("%s: %w", comp.Name(), err))
		} else {
			l.logger.Info("component stopped", "name", comp.Name())
		}
	}

	return errs
}

// indexOf returns the index of a component in a slice
func indexOf(components []Component, target Component) int {
	for i, comp := range components {
		if comp == target {
			return i
		}
	}
	return -1
}

// Runtime manages the UMC runtime with signal handling
type Runtime struct {
	launcher *Launcher
	logger   *slog.Logger
}

// NewRuntime creates a new runtime manager
func NewRuntime(launcher *Launcher) *Runtime {
	return &Runtime{
		launcher: launcher,
		logger:   launcher.logger,
	}
}

// Run starts the launcher and blocks until a termination signal is received
func (r *Runtime) Run() error {
	// Start all components
	ctx := context.Background()
	if err := r.launcher.Start(ctx); err != nil {
		return err
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	r.logger.Info("runtime started, waiting for signals")

	for {
		sig := <-sigChan
		r.logger.Info("received signal", "signal", sig)

		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			r.logger.Info("shutting down")
			if err := r.launcher.Stop(ctx); err != nil {
				r.logger.Error("shutdown error", "error", err)
				return err
			}
			return nil
		case syscall.SIGHUP:
			// UMCs can implement config reload in their components
			// This is a signal to trigger reload, not automatic
			r.logger.Info("received SIGHUP, reloading configuration")
		}
	}
}

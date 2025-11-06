package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/dexguitar/spacecraftory/platform/pkg/logger"
)

// shutdownTimeout default value, can be made a parameter
const shutdownTimeout = 5 * time.Second

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

// Closer manages the graceful shutdown process of the application
type Closer struct {
	mu     sync.Mutex                    // Protection against race conditions when adding functions
	once   sync.Once                     // Guarantees single invocation of CloseAll
	done   chan struct{}                 // Channel for completion notification
	funcs  []func(context.Context) error // Registered close functions
	logger Logger                        // Logger being used
}

// Global instance for use throughout the application
var globalCloser = NewWithLogger(&logger.NoopLogger{})

// AddNamed adds a close function with a dependency name for logging to the global closer
func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

// Add adds close functions to the global closer
func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

// CloseAll initiates the closing process of all registered functions in the global closer
func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

// SetLogger allows setting a custom logger for the global closer
func SetLogger(l Logger) {
	globalCloser.SetLogger(l)
}

// Configure sets up the global closer to handle system signals
func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

// New creates a new Closer instance with the default logger
func New(signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Logger(), signals...)
}

// NewWithLogger creates a new Closer instance with a specified logger.
// If signals are provided, Closer will start listening to them and call CloseAll upon receipt.
func NewWithLogger(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// SetLogger sets the logger for Closer
func (c *Closer) SetLogger(l Logger) {
	c.logger = l
}

// handleSignals handles system signals and calls CloseAll with fresh shutdown context
func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		c.logger.Info(context.Background(), "🛑 System signal received, starting graceful shutdown...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "❌ Error closing resources: %v", zap.Error(err))
		}

	case <-c.done:
		// CloseAll was already called manually, just exit
	}
}

// AddNamed adds a close function with a dependency name for logging
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(func(ctx context.Context) error {
		start := time.Now()
		c.logger.Info(ctx, fmt.Sprintf("🧩 Closing %s...", name))

		err := f(ctx)

		duration := time.Since(start)
		if err != nil {
			c.logger.Error(ctx, fmt.Sprintf("❌ Error closing %s: %v (took %s)", name, err, duration))
		} else {
			c.logger.Info(ctx, fmt.Sprintf("✅ %s successfully closed in %s", name, duration))
		}
		return err
	})
}

// Add adds one or more close functions
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// CloseAll calls all registered close functions.
// Returns the first error encountered, if any.
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil // free up memory
		c.mu.Unlock()

		if len(funcs) == 0 {
			c.logger.Info(ctx, "ℹ️ No functions to close.")
			return
		}

		c.logger.Info(ctx, "🚦 Starting graceful shutdown process...")

		errCh := make(chan error, len(funcs))
		var wg sync.WaitGroup

		// Execute in reverse order of addition
		for i := len(funcs) - 1; i >= 0; i-- {
			f := funcs[i]
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()

				// Panic protection
				defer func() {
					if r := recover(); r != nil {
						errCh <- errors.New("panic recovered in closer")
						c.logger.Error(ctx, "⚠️ Panic in close function", zap.Any("error", r))
					}
				}()

				if err := f(ctx); err != nil {
					errCh <- err
				}
			}(f)
		}

		// Close error channel when all functions complete
		go func() {
			wg.Wait()
			close(errCh)
		}()

		// Read errors or context cancellation
		for {
			select {
			case <-ctx.Done():
				c.logger.Info(ctx, "⚠️ Context cancelled during closing", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, "✅ All resources successfully closed")
					return
				}
				c.logger.Error(ctx, "❌ Error during closing", zap.Error(err))
				if result == nil {
					result = err
				}
			}
		}
	})

	return result
}

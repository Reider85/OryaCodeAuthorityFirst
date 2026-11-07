package ui

import (
	"github.com/charmbracelet/bubbletea"
	"time"
)

// PerformanceMiddleware provides a middleware for performance monitoring
type PerformanceMiddleware struct {
	monitor *PerformanceMonitor
}

// NewPerformanceMiddleware creates a new performance middleware
func NewPerformanceMiddleware(monitor *PerformanceMonitor) *PerformanceMiddleware {
	return &PerformanceMiddleware{
		monitor: monitor,
	}
}

// MeasureOperation measures an operation and returns the result
func (pm *PerformanceMiddleware) MeasureOperation(name string, fn func() interface{}) interface{} {
	ctx := pm.monitor.StartOperation(name)
	result := fn()
	ctx.End(nil)
	return result
}

// MeasureOperationWithDuration measures an operation and returns both result and duration
func (pm *PerformanceMiddleware) MeasureOperationWithDuration(name string, fn func() interface{}) (interface{}, time.Duration) {
	ctx := pm.monitor.StartOperation(name)
	result := fn()
	duration := time.Since(ctx.startTime)
	ctx.End(nil)
	return result, duration
}

// MeasureAsyncOperation measures an async operation
func (pm *PerformanceMiddleware) MeasureAsyncOperation(name string, fn func()) {
	go func() {
		ctx := pm.monitor.StartOperation(name)
		fn()
		ctx.End(nil)
	}()
}

// TeaCmd creates a tea.Cmd for measuring performance
func (pm *PerformanceMiddleware) TeaCmd(name string, fn func() tea.Msg) tea.Cmd {
	return func() tea.Msg {
		ctx := pm.monitor.StartOperation(name)
		msg := fn()
		ctx.End(nil)
		return msg
	}
}

// TeaCmdWithMetadata creates a tea.Cmd for measuring performance with metadata
func (pm *PerformanceMiddleware) TeaCmdWithMetadata(name string, metadata map[string]interface{}, fn func() tea.Msg) tea.Cmd {
	return func() tea.Msg {
		ctx := pm.monitor.StartOperation(name)
		msg := fn()
		ctx.End(metadata)
		return msg
	}
}

// TeaCmdWithTimeout creates a tea.Cmd for measuring performance with timeout
func (pm *PerformanceMiddleware) TeaCmdWithTimeout(name string, timeout time.Duration, fn func() tea.Msg) tea.Cmd {
	return func() tea.Msg {
		ctx := pm.monitor.StartOperation(name)
		
		// Create a channel to receive the result
		resultCh := make(chan tea.Msg)
		go func() {
			resultCh <- fn()
		}()
		
		// Wait for result or timeout
		select {
		case msg := <-resultCh:
			ctx.End(nil)
			return msg
		case <-time.After(timeout):
			ctx.End(map[string]interface{}{
				"timeout": true,
				"timeout_duration": timeout,
			})
			return tea.QuitMsg{}
		}
	}
}

// BatchMeasure measures multiple operations and returns their results
func (pm *PerformanceMiddleware) BatchMeasure(operations []struct {
	Name string
	Fn   func() interface{}
}) []interface{} {
	results := make([]interface{}, len(operations))
	
	for i, op := range operations {
		results[i] = pm.MeasureOperation(op.Name, op.Fn)
	}
	
	return results
}

// BatchMeasureWithDuration measures multiple operations and returns their results and durations
func (pm *PerformanceMiddleware) BatchMeasureWithDuration(operations []struct {
	Name string
	Fn   func() interface{}
}) ([]interface{}, []time.Duration) {
	results := make([]interface{}, len(operations))
	durations := make([]time.Duration, len(operations))
	
	for i, op := range operations {
		result, duration := pm.MeasureOperationWithDuration(op.Name, op.Fn)
		results[i] = result
		durations[i] = duration
	}
	
	return results, durations
}

// ConcurrentMeasure measures multiple operations concurrently
func (pm *PerformanceMiddleware) ConcurrentMeasure(operations []struct {
	Name string
	Fn   func() interface{}
}) []interface{} {
	results := make([]interface{}, len(operations))
	
	// Create a channel to collect results
	resultCh := make(chan struct {
		index int
		result interface{}
	}, len(operations))
	
	// Start all operations concurrently
	for i, op := range operations {
		go func(idx int, name string, fn func() interface{}) {
			result := pm.MeasureOperation(name, fn)
			resultCh <- struct {
				index int
				result interface{}
			}{idx, result}
		}(i, op.Name, op.Fn)
	}
	
	// Collect results
	for i := 0; i < len(operations); i++ {
		res := <-resultCh
		results[res.index] = res.result
	}
	
	return results
}

// ConditionalMeasure measures an operation only if condition is met
func (pm *PerformanceMiddleware) ConditionalMeasure(name string, condition func() bool, fn func() interface{}) interface{} {
	if !condition() {
		return nil
	}
	
	ctx := pm.monitor.StartOperation(name)
	result := fn()
	ctx.End(nil)
	return result
}

// RetryMeasure measures an operation with retry logic
func (pm *PerformanceMiddleware) RetryMeasure(name string, maxRetries int, fn func() interface{}) interface{} {
	var result interface{}
	var lastError error
	
	for i := 0; i < maxRetries; i++ {
		ctx := pm.monitor.StartOperation(name)
		result = fn()
		ctx.End(nil)
		
		// Check if operation was successful
		if result != nil {
			return result
		}
		
		// Add delay between retries
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}
	
	return result
}

// CircuitBreakerMeasure measures an operation with circuit breaker pattern
func (pm *PerformanceMiddleware) CircuitBreakerMeasure(name string, failureThreshold int, fn func() interface{}) interface{} {
	// Simple circuit breaker implementation
	failures := 0
	
	// Check if circuit is broken
	if failures >= failureThreshold {
		// Return cached result or error
		return nil
	}
	
	ctx := pm.monitor.StartOperation(name)
	result := fn()
	ctx.End(nil)
	
	// Check if operation failed
	if result == nil {
		failures++
	} else {
		failures = 0 // Reset on success
	}
	
	return result
}
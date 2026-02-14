package logger

import (
	"log/slog"
	"sync/atomic"
)

var (
	// samplingRate is the current sampling rate (0.0-1.0)
	samplingRate atomic.Value // stores float64
	// samplingCounter is used for deterministic sampling
	samplingCounter atomic.Uint64
)

// setSamplingRate sets the global sampling rate
func setSamplingRate(rate float64) {
	// Clamp rate between 0.0 and 1.0
	if rate < 0.0 {
		rate = 0.0
	} else if rate > 1.0 {
		rate = 1.0
	}
	samplingRate.Store(rate)
}

// getSamplingRate gets the current sampling rate
func getSamplingRate() float64 {
	rate := samplingRate.Load()
	if rate == nil {
		return 1.0 // Default to log everything
	}
	return rate.(float64)
}

// shouldSample determines if a log entry should be logged based on sampling rate
// Always logs errors and fatal messages regardless of sampling rate
func shouldSample(level slog.Level, rate float64) bool {
	// Always log errors and above
	if level >= slog.LevelError {
		return true
	}

	// If rate is 1.0, log everything
	if rate >= 1.0 {
		return true
	}

	// If rate is 0.0, log nothing (except errors handled above)
	if rate <= 0.0 {
		return false
	}

	// Use deterministic sampling based on counter for consistency
	counter := samplingCounter.Add(1)
	// Use modulo for deterministic sampling
	threshold := uint64(1.0 / rate)
	return (counter % threshold) == 0
}

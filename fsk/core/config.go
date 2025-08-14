// Package core implements the pure FSK algorithm without any external dependencies.
// This package contains the core signal processing logic that can be used
// across different platforms including WebAssembly.
package core

// Config holds the FSK modem configuration parameters.
type Config struct {
	BaseFreq    float64 // Base frequency in Hz
	FreqSpacing float64 // Frequency spacing between symbols
	Order       int     // FSK order (2^n symbols)
	BaudRate    float64 // Symbol rate (symbols per second)
	SampleRate  int     // Audio sample rate
}

// DefaultConfig returns a default FSK configuration.
func DefaultConfig() Config {
	return Config{
		BaseFreq:    1000,
		FreqSpacing: 200,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}
}

// UltrasonicConfig returns a configuration optimized for ultrasonic communication.
func UltrasonicConfig() Config {
	return Config{
		BaseFreq:    22000,
		FreqSpacing: 500,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}
}
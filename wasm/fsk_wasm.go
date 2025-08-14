// WASM-compatible FSK implementation (encoding only)
package main

import (
	"math"
)

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

// Modem represents an FSK modem with encoding capabilities.
type Modem struct {
	config       Config
	symbolPeriod int // Samples per symbol
	frequencies  []float64
	phase        []float64 // Phase accumulators for each frequency
}

// New creates a new FSK modem with the given configuration.
func New(config Config) *Modem {
	modem := &Modem{
		config:       config,
		symbolPeriod: int(float64(config.SampleRate) / config.BaudRate),
		frequencies:  make([]float64, 1<<config.Order), // 2^order frequencies
		phase:        make([]float64, 1<<config.Order),
	}

	// Calculate frequencies for each symbol
	for i := 0; i < len(modem.frequencies); i++ {
		modem.frequencies[i] = config.BaseFreq + float64(i)*config.FreqSpacing
	}

	return modem
}

// Config returns the modem's configuration.
func (m *Modem) Config() Config {
	return m.config
}

// Frequencies returns the array of frequencies used by this modem.
func (m *Modem) Frequencies() []float64 {
	return append([]float64(nil), m.frequencies...) // Return copy
}

// SymbolPeriod returns the number of samples per symbol.
func (m *Modem) SymbolPeriod() int {
	return m.symbolPeriod
}

// Encode converts binary data to FSK-modulated audio signal.
func (m *Modem) Encode(data []byte) []float32 {
	bitsPerSymbol := m.config.Order
	totalBits := len(data) * 8
	symbolCount := (totalBits + bitsPerSymbol - 1) / bitsPerSymbol // Ceiling division

	output := make([]float32, symbolCount*m.symbolPeriod)

	bitIndex := 0
	for symbolIdx := 0; symbolIdx < symbolCount; symbolIdx++ {
		// Extract bits for this symbol
		symbol := 0
		for bit := 0; bit < bitsPerSymbol && bitIndex < totalBits; bit++ {
			byteIdx := bitIndex / 8
			bitInByte := 7 - (bitIndex % 8) // MSB first

			if data[byteIdx]&(1<<bitInByte) != 0 {
				symbol |= 1 << (bitsPerSymbol - 1 - bit)
			}
			bitIndex++
		}

		// Generate waveform for this symbol
		freq := m.frequencies[symbol]
		phaseIncrement := 2 * math.Pi * freq / float64(m.config.SampleRate)

		for sampleIdx := 0; sampleIdx < m.symbolPeriod; sampleIdx++ {
			outputIdx := symbolIdx*m.symbolPeriod + sampleIdx
			if outputIdx < len(output) {
				output[outputIdx] = float32(0.5 * math.Sin(m.phase[symbol]))
				m.phase[symbol] += phaseIncrement

				// Keep phase in range [0, 2π]
				if m.phase[symbol] >= 2*math.Pi {
					m.phase[symbol] -= 2 * math.Pi
				}
			}
		}
	}

	return output
}
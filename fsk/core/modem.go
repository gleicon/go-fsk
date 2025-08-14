package core

import "math"

// Modem represents an FSK modem with encoding/decoding capabilities.
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

// correlateWithFrequency calculates correlation between signal and reference frequency.
func (m *Modem) correlateWithFrequency(signal []float32, freq float64) float64 {
	phaseIncrement := 2 * math.Pi * freq / float64(m.config.SampleRate)

	var correlation float64
	phase := 0.0

	for _, sample := range signal {
		reference := math.Sin(phase)
		correlation += float64(sample) * reference
		phase += phaseIncrement

		if phase >= 2*math.Pi {
			phase -= 2 * math.Pi
		}
	}

	return math.Abs(correlation) / float64(len(signal))
}
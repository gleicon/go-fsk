package core

import "math"

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
package core

// Decode converts FSK-modulated audio signal back to binary data.
func (m *Modem) Decode(signal []float32) []byte {
	symbolCount := len(signal) / m.symbolPeriod
	if symbolCount == 0 {
		return nil
	}

	symbols := make([]int, symbolCount)

	// For each symbol period, determine which frequency has the highest correlation
	for symbolIdx := 0; symbolIdx < symbolCount; symbolIdx++ {
		start := symbolIdx * m.symbolPeriod
		end := start + m.symbolPeriod
		if end > len(signal) {
			end = len(signal)
		}

		maxCorrelation := -1.0
		detectedSymbol := 0

		// Test correlation with each possible frequency
		for freqIdx, freq := range m.frequencies {
			correlation := m.correlateWithFrequency(signal[start:end], freq)
			if correlation > maxCorrelation {
				maxCorrelation = correlation
				detectedSymbol = freqIdx
			}
		}

		symbols[symbolIdx] = detectedSymbol
	}

	// Convert symbols back to bytes
	bitsPerSymbol := m.config.Order
	totalBits := symbolCount * bitsPerSymbol
	byteCount := (totalBits + 7) / 8 // Ceiling division

	output := make([]byte, byteCount)

	bitIndex := 0
	for _, symbol := range symbols {
		for bit := 0; bit < bitsPerSymbol && bitIndex < totalBits; bit++ {
			byteIdx := bitIndex / 8
			bitInByte := 7 - (bitIndex % 8) // MSB first

			if symbol&(1<<(bitsPerSymbol-1-bit)) != 0 {
				output[byteIdx] |= 1 << bitInByte
			}
			bitIndex++
		}
	}

	return output
}
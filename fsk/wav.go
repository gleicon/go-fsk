package fsk

import (
	"encoding/binary"
	"fmt"
	"os"
)

// WriteWAVFile writes audio samples to a WAV file (16-bit PCM).
func (m *Modem) WriteWAVFile(filename string, signal []float32) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// WAV header
	header := []byte("RIFF")
	file.Write(header)

	// File size (will update later)
	dataSize := len(signal) * 2 // 16-bit samples
	fileSize := uint32(36 + dataSize)
	binary.Write(file, binary.LittleEndian, fileSize)

	// Format
	file.Write([]byte("WAVE"))
	file.Write([]byte("fmt "))

	// Format chunk size
	binary.Write(file, binary.LittleEndian, uint32(16))

	// Audio format (1 = PCM)
	binary.Write(file, binary.LittleEndian, uint16(1))

	// Channels
	binary.Write(file, binary.LittleEndian, uint16(1))

	// Sample rate
	binary.Write(file, binary.LittleEndian, uint32(m.config.SampleRate))

	// Byte rate
	byteRate := uint32(m.config.SampleRate * 2)
	binary.Write(file, binary.LittleEndian, byteRate)

	// Block align
	binary.Write(file, binary.LittleEndian, uint16(2))

	// Bits per sample
	binary.Write(file, binary.LittleEndian, uint16(16))

	// Data chunk
	file.Write([]byte("data"))
	binary.Write(file, binary.LittleEndian, uint32(dataSize))

	// Write audio data (convert float32 to int16)
	for _, sample := range signal {
		// Clamp and convert to 16-bit
		if sample > 1.0 {
			sample = 1.0
		}
		if sample < -1.0 {
			sample = -1.0
		}
		intSample := int16(sample * 32767)
		binary.Write(file, binary.LittleEndian, intSample)
	}

	return nil
}

// ReadWAVFile reads audio samples from a WAV file (16-bit PCM).
func (m *Modem) ReadWAVFile(filename string) ([]float32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Skip WAV header (44 bytes for standard PCM)
	header := make([]byte, 44)
	_, err = file.Read(header)
	if err != nil {
		return nil, err
	}

	// Verify it's a WAV file
	if string(header[0:4]) != "RIFF" || string(header[8:12]) != "WAVE" {
		return nil, fmt.Errorf("not a valid WAV file")
	}

	// Read audio data
	var signal []float32
	for {
		var sample int16
		err := binary.Read(file, binary.LittleEndian, &sample)
		if err != nil {
			break // End of file
		}
		// Convert to float32
		floatSample := float32(sample) / 32767.0
		signal = append(signal, floatSample)
	}

	return signal, nil
}

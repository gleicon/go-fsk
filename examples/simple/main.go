// Simple FSK encoding/decoding example
package main

import (
	"fmt"
	"log"

	"github.com/gleicon/go-fsk/fsk/core"
	"github.com/gleicon/go-fsk/fsk/utils"
)

func main() {
	// Create FSK modem with default configuration
	config := core.DefaultConfig()
	modem := core.New(config)

	// Display configuration
	fmt.Printf("FSK Configuration:\n")
	fmt.Printf("  Base Frequency: %.0f Hz\n", config.BaseFreq)
	fmt.Printf("  Frequency Spacing: %.0f Hz\n", config.FreqSpacing)
	fmt.Printf("  FSK Order: %d (2^%d = %d symbols)\n", config.Order, config.Order, 1<<config.Order)
	fmt.Printf("  Baud Rate: %.0f symbols/sec\n", config.BaudRate)
	fmt.Printf("  Sample Rate: %d Hz\n", config.SampleRate)
	fmt.Printf("\n")

	// Test message
	message := "Hello, FSK World! 🎵"
	fmt.Printf("Original message: %s\n", message)

	// Encode message to FSK signal
	signal := modem.Encode([]byte(message))
	fmt.Printf("Encoded to %d audio samples\n", len(signal))

	// Decode signal back to message
	decoded := modem.Decode(signal)
	decodedMessage := string(decoded)
	fmt.Printf("Decoded message:  %s\n", decodedMessage)

	// Verify accuracy
	if decodedMessage == message {
		fmt.Println("✅ Encoding/Decoding successful!")
	} else {
		fmt.Println("❌ Encoding/Decoding failed!")
		log.Printf("Expected: %q", message)
		log.Printf("Got:      %q", decodedMessage)
	}

	// Save to WAV file for testing
	fmt.Printf("\nSaving signal to test.wav...\n")
	err := utils.WriteWAVFile("test.wav", signal, modem.Config())
	if err != nil {
		log.Printf("Error writing WAV file: %v", err)
	} else {
		fmt.Printf("Signal saved to test.wav (play with: ffplay test.wav)\n")
	}
}

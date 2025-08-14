// Ultrasonic FSK communication example
package main

import (
	"fmt"
	"time"

	"github.com/gleicon/go-fsk/fsk/core"
	"github.com/gleicon/go-fsk/fsk/realtime"
	"github.com/gleicon/go-fsk/fsk/utils"
)

func main() {
	// Configure for ultrasonic communication (22kHz+)
	config := core.UltrasonicConfig()
	config.BaseFreq = 22000  // 22kHz base (inaudible to most humans)
	config.FreqSpacing = 500 // 500Hz spacing
	config.BaudRate = 100    // 100 symbols/second

	modem := core.New(config)
	frequencies := modem.Frequencies()

	fmt.Printf("Ultrasonic FSK Modem\n")
	fmt.Printf("====================\n")
	fmt.Printf("Base Frequency: %.0f Hz\n", config.BaseFreq)
	fmt.Printf("Frequencies: ")
	for i, freq := range frequencies {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%.0f", freq)
	}
	fmt.Printf(" Hz\n")
	fmt.Printf("Human audible: %t\n", frequencies[0] < 20000)
	fmt.Printf("\n")

	// Test messages
	messages := []string{
		"Secret message 1",
		"Steganographic data",
		"Ultrasonic beacon",
		"Covert channel active",
	}

	for i, message := range messages {
		fmt.Printf("Test %d: %s\n", i+1, message)

		// Encode message
		signal := modem.Encode([]byte(message))
		duration := float64(len(signal)) / float64(config.SampleRate)
		fmt.Printf("  Encoded to %d samples (%.2f seconds)\n", len(signal), duration)

		// Decode message
		decoded := modem.Decode(signal)
		decodedMessage := string(decoded)

		// Verify
		if decodedMessage == message {
			fmt.Printf("  ✅ Decoded successfully: %s\n", decodedMessage)
		} else {
			fmt.Printf("  ❌ Decode failed: got %q\n", decodedMessage)
		}

		// Save each test as WAV file
		filename := fmt.Sprintf("ultrasonic_test_%d.wav", i+1)
		err := utils.WriteWAVFile(filename, signal, modem.Config())
		if err != nil {
			fmt.Printf("  Error saving %s: %v\n", filename, err)
		} else {
			fmt.Printf("  Saved to %s\n", filename)
		}

		fmt.Println()
	}

	// Demonstrate real-time ultrasonic transmission
	fmt.Printf("Real-time ultrasonic transmission test:\n")
	fmt.Printf("=======================================\n")

	// Create transmitter
	transmitter, err := realtime.NewTransmitter(modem)
	if err != nil {
		fmt.Printf("Failed to create transmitter: %v\n", err)
		return
	}
	defer transmitter.Close()

	testMessage := "Ultrasonic test signal"
	fmt.Printf("Transmitting: %s\n", testMessage)
	fmt.Printf("Frequency range: %.0f-%.0f Hz (ultrasonic)\n",
		frequencies[0], frequencies[len(frequencies)-1])

	// Transmit (will be inaudible at 22kHz+)
	err = transmitter.Transmit([]byte(testMessage))
	if err != nil {
		fmt.Printf("Transmission failed: %v\n", err)
	} else {
		fmt.Printf("✅ Ultrasonic transmission complete!\n")
		fmt.Printf("(Signal was transmitted at %dkHz - inaudible to humans)\n", int(config.BaseFreq/1000))
	}

	// Wait a moment
	time.Sleep(1 * time.Second)

	fmt.Printf("\nUltrasonic FSK can be used for:\n")
	fmt.Printf("- Steganographic communication\n")
	fmt.Printf("- Device-to-device pairing\n")
	fmt.Printf("- Covert data channels\n")
	fmt.Printf("- Audio watermarking\n")
	fmt.Printf("- IoT sensor networks\n")
}

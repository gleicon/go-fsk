// Frequency collision and mixing test
package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/gleicon/go-fsk/fsk"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <scenario>")
		fmt.Println("Scenarios:")
		fmt.Println("  1: Same frequency collision")
		fmt.Println("  2: Overlapping frequencies")
		fmt.Println("  3: Separate frequencies (clean)")
		fmt.Println("  4: Multi-channel broadcast")
		fmt.Println("  5: Point-to-point duplex")
		return
	}

	scenario, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid scenario: %v\n", err)
		return
	}

	switch scenario {
	case 1:
		testSameFrequencyCollision()
	case 2:
		testOverlappingFrequencies()
	case 3:
		testSeparateFrequencies()
	case 4:
		testMultiChannelBroadcast()
	case 5:
		testPointToPointDuplex()
	default:
		fmt.Printf("Unknown scenario: %d\n", scenario)
	}
}

func testSameFrequencyCollision() {
	fmt.Println("=== Same Frequency Collision Test ===")
	fmt.Println("Two modems using identical frequencies - expect collisions")

	// Both modems use same frequency
	config := fsk.Config{
		BaseFreq:    22000,
		FreqSpacing: 500,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}

	modem1 := fsk.New(config)
	modem2 := fsk.New(config)

	message1 := "Message from Agent A"
	message2 := "Message from Agent B"

	// Encode both messages
	signal1 := modem1.Encode([]byte(message1))
	signal2 := modem2.Encode([]byte(message2))

	fmt.Printf("Agent A: %s\n", message1)
	fmt.Printf("Agent B: %s\n", message2)
	fmt.Printf("Both using frequencies: %.0f-%.0f Hz\n",
		config.BaseFreq, config.BaseFreq+3*config.FreqSpacing)

	// Mix signals (collision)
	mixed := make([]float32, len(signal1))
	for i := 0; i < len(mixed) && i < len(signal2); i++ {
		mixed[i] = signal1[i] + signal2[i] // Direct collision
	}

	// Try to decode mixed signal
	decoded1 := modem1.Decode(mixed)
	decoded2 := modem2.Decode(mixed)

	fmt.Printf("\nAfter collision:\n")
	fmt.Printf("Modem 1 decoded: %q\n", string(decoded1))
	fmt.Printf("Modem 2 decoded: %q\n", string(decoded2))

	// Check corruption
	success1 := string(decoded1) == message1
	success2 := string(decoded2) == message2

	fmt.Printf("\nResults:\n")
	fmt.Printf("Agent A recovery: %t\n", success1)
	fmt.Printf("Agent B recovery: %t\n", success2)

	if !success1 && !success2 {
		fmt.Println("COLLISION: Both messages corrupted")
	}

	// Save mixed signal for analysis
	err := modem1.WriteWAVFile("collision_mixed.wav", mixed)
	if err == nil {
		fmt.Println("Mixed signal saved to collision_mixed.wav")
	}
}

func testOverlappingFrequencies() {
	fmt.Println("=== Overlapping Frequencies Test ===")
	fmt.Println("Two modems with overlapping frequency ranges")

	// Overlapping frequency ranges
	config1 := fsk.Config{
		BaseFreq:    22000, // 22000-23000 Hz
		FreqSpacing: 300,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}

	config2 := fsk.Config{
		BaseFreq:    22500, // 22500-23500 Hz (overlaps)
		FreqSpacing: 300,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}

	modem1 := fsk.New(config1)
	modem2 := fsk.New(config2)

	message1 := "Agent A message"
	message2 := "Agent B message"

	signal1 := modem1.Encode([]byte(message1))
	signal2 := modem2.Encode([]byte(message2))

	fmt.Printf("Agent A: %s (%.0f-%.0f Hz)\n", message1,
		config1.BaseFreq, config1.BaseFreq+3*config1.FreqSpacing)
	fmt.Printf("Agent B: %s (%.0f-%.0f Hz)\n", message2,
		config2.BaseFreq, config2.BaseFreq+3*config2.FreqSpacing)
	fmt.Printf("Overlap: %.0f-%.0f Hz\n",
		math.Max(config1.BaseFreq, config2.BaseFreq),
		math.Min(config1.BaseFreq+3*config1.FreqSpacing, config2.BaseFreq+3*config2.FreqSpacing))

	// Mix signals with some interference
	mixed := make([]float32, len(signal1))
	for i := 0; i < len(mixed) && i < len(signal2); i++ {
		mixed[i] = signal1[i] + 0.7*signal2[i] // Partial interference
	}

	decoded1 := modem1.Decode(mixed)
	decoded2 := modem2.Decode(mixed)

	fmt.Printf("\nAfter partial interference:\n")
	fmt.Printf("Modem 1 decoded: %q\n", string(decoded1))
	fmt.Printf("Modem 2 decoded: %q\n", string(decoded2))

	success1 := string(decoded1) == message1
	success2 := string(decoded2) == message2

	fmt.Printf("\nResults:\n")
	fmt.Printf("Agent A recovery: %t\n", success1)
	fmt.Printf("Agent B recovery: %t\n", success2)

	err := modem1.WriteWAVFile("overlap_mixed.wav", mixed)
	if err == nil {
		fmt.Println("Mixed signal saved to overlap_mixed.wav")
	}
}

func testSeparateFrequencies() {
	fmt.Println("=== Separate Frequencies Test ===")
	fmt.Println("Two modems using well-separated frequency ranges")

	config1 := fsk.Config{
		BaseFreq:    22000, // 22000-23000 Hz
		FreqSpacing: 250,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}

	config2 := fsk.Config{
		BaseFreq:    24000, // 24000-25000 Hz (2kHz separation)
		FreqSpacing: 250,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}

	modem1 := fsk.New(config1)
	modem2 := fsk.New(config2)

	message1 := "Clean message from A"
	message2 := "Clean message from B"

	signal1 := modem1.Encode([]byte(message1))
	signal2 := modem2.Encode([]byte(message2))

	fmt.Printf("Agent A: %s (%.0f-%.0f Hz)\n", message1,
		config1.BaseFreq, config1.BaseFreq+3*config1.FreqSpacing)
	fmt.Printf("Agent B: %s (%.0f-%.0f Hz)\n", message2,
		config2.BaseFreq, config2.BaseFreq+3*config2.FreqSpacing)
	fmt.Printf("Separation: %.0f Hz\n",
		config2.BaseFreq-config1.BaseFreq-3*config1.FreqSpacing)

	// Mix signals - should not interfere
	mixed := make([]float32, len(signal1))
	for i := 0; i < len(mixed) && i < len(signal2); i++ {
		mixed[i] = signal1[i] + signal2[i]
	}

	// Each modem should decode its own signal cleanly
	decoded1 := modem1.Decode(mixed)
	decoded2 := modem2.Decode(mixed)

	fmt.Printf("\nWith both signals present:\n")
	fmt.Printf("Modem 1 decoded: %q\n", string(decoded1))
	fmt.Printf("Modem 2 decoded: %q\n", string(decoded2))

	success1 := string(decoded1) == message1
	success2 := string(decoded2) == message2

	fmt.Printf("\nResults:\n")
	fmt.Printf("Agent A recovery: %t\n", success1)
	fmt.Printf("Agent B recovery: %t\n", success2)

	if success1 && success2 {
		fmt.Println("SUCCESS: Clean separation allows both signals")
	}

	err := modem1.WriteWAVFile("clean_mixed.wav", mixed)
	if err == nil {
		fmt.Println("Mixed signal saved to clean_mixed.wav")
	}
}

func testMultiChannelBroadcast() {
	fmt.Println("=== Multi-Channel Broadcast Test ===")
	fmt.Println("Multiple channels with different users")

	channels := []fsk.Config{
		{BaseFreq: 22000, FreqSpacing: 500, Order: 2, BaudRate: 100, SampleRate: 48000},
		{BaseFreq: 24000, FreqSpacing: 500, Order: 2, BaudRate: 100, SampleRate: 48000},
		{BaseFreq: 26000, FreqSpacing: 500, Order: 2, BaudRate: 100, SampleRate: 48000},
	}

	messages := []string{
		"Channel 1 broadcast",
		"Channel 2 broadcast",
		"Channel 3 broadcast",
	}

	var signals [][]float32
	var modems []*fsk.Modem

	for i, config := range channels {
		modem := fsk.New(config)
		signal := modem.Encode([]byte(messages[i]))
		modems = append(modems, modem)
		signals = append(signals, signal)

		fmt.Printf("Channel %d: %s (%.0fkHz)\n",
			i+1, messages[i], config.BaseFreq/1000)
	}

	// Mix all signals
	maxLen := 0
	for _, signal := range signals {
		if len(signal) > maxLen {
			maxLen = len(signal)
		}
	}

	mixed := make([]float32, maxLen)
	for _, signal := range signals {
		for i, sample := range signal {
			mixed[i] += sample
		}
	}

	fmt.Println("\nDecoding each channel from mixed signal:")
	for i, modem := range modems {
		decoded := modem.Decode(mixed)
		success := string(decoded) == messages[i]
		fmt.Printf("Channel %d: %q (success: %t)\n",
			i+1, string(decoded), success)
	}

	err := modems[0].WriteWAVFile("multichannel_broadcast.wav", mixed)
	if err == nil {
		fmt.Println("Multi-channel signal saved to multichannel_broadcast.wav")
	}
}

func testPointToPointDuplex() {
	fmt.Println("=== Point-to-Point Duplex Test ===")
	fmt.Println("Two agents with dedicated TX/RX frequency pairs")

	// Agent A: TX on 22kHz, RX on 24kHz
	// Agent B: TX on 24kHz, RX on 22kHz
	configA_TX := fsk.Config{BaseFreq: 22000, FreqSpacing: 500, Order: 2, BaudRate: 100, SampleRate: 48000}
	configA_RX := fsk.Config{BaseFreq: 24000, FreqSpacing: 500, Order: 2, BaudRate: 100, SampleRate: 48000}

	configB_TX := fsk.Config{BaseFreq: 24000, FreqSpacing: 500, Order: 2, BaudRate: 100, SampleRate: 48000}
	configB_RX := fsk.Config{BaseFreq: 22000, FreqSpacing: 500, Order: 2, BaudRate: 100, SampleRate: 48000}

	modemA_TX := fsk.New(configA_TX)
	modemA_RX := fsk.New(configA_RX)
	modemB_TX := fsk.New(configB_TX)
	modemB_RX := fsk.New(configB_RX)

	messageA_to_B := "A to B: Hello"
	messageB_to_A := "B to A: Hi there"

	signalA := modemA_TX.Encode([]byte(messageA_to_B))
	signalB := modemB_TX.Encode([]byte(messageB_to_A))

	fmt.Printf("Agent A TX: %s (%.0fkHz)\n", messageA_to_B, configA_TX.BaseFreq/1000)
	fmt.Printf("Agent B TX: %s (%.0fkHz)\n", messageB_to_A, configB_TX.BaseFreq/1000)

	// Simultaneous transmission
	maxLen := len(signalA)
	if len(signalB) > maxLen {
		maxLen = len(signalB)
	}

	mixed := make([]float32, maxLen)
	for i := 0; i < maxLen; i++ {
		if i < len(signalA) {
			mixed[i] += signalA[i]
		}
		if i < len(signalB) {
			mixed[i] += signalB[i]
		}
	}

	// Each agent receives on their RX frequency
	decodedAt_A := modemA_RX.Decode(mixed) // A receives B's message
	decodedAt_B := modemB_RX.Decode(mixed) // B receives A's message

	fmt.Printf("\nDuplex communication results:\n")
	fmt.Printf("Agent A received: %q\n", string(decodedAt_A))
	fmt.Printf("Agent B received: %q\n", string(decodedAt_B))

	successA := string(decodedAt_A) == messageB_to_A
	successB := string(decodedAt_B) == messageA_to_B

	fmt.Printf("\nResults:\n")
	fmt.Printf("A->B communication: %t\n", successB)
	fmt.Printf("B->A communication: %t\n", successA)

	if successA && successB {
		fmt.Println("SUCCESS: Full duplex communication works!")
	}

	err := modemA_TX.WriteWAVFile("duplex_mixed.wav", mixed)
	if err == nil {
		fmt.Println("Duplex signal saved to duplex_mixed.wav")
	}

	// Demonstrate why this works
	fmt.Printf("\nFrequency Analysis:\n")
	fmt.Printf("A transmits %.0fkHz, B listens %.0fkHz: %t\n",
		configA_TX.BaseFreq/1000, configB_RX.BaseFreq/1000,
		configA_TX.BaseFreq == configB_RX.BaseFreq)
	fmt.Printf("B transmits %.0fkHz, A listens %.0fkHz: %t\n",
		configB_TX.BaseFreq/1000, configA_RX.BaseFreq/1000,
		configB_TX.BaseFreq == configA_RX.BaseFreq)
}

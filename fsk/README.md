# FSK Package

A high-performance Go package for Frequency Shift Keying (FSK) modulation and demodulation. Supports both audible and ultrasonic frequencies with real-time audio capabilities.

## Features

- **High-order FSK**: Configurable 2^n symbols (2, 4, 8, 16 symbols)
- **Frequency agnostic**: Works from audible (1kHz) to ultrasonic (25kHz+) 
- **Real-time audio**: Live microphone capture and speaker playback
- **WAV file support**: Generate and decode audio files
- **Duplex communication**: Full-duplex chat mode
- **Cross-platform**: Windows, macOS, Linux support via Malgo

## Installation

```bash
go get github.com/gleicon/go-fsk/fsk
```

## Quick Start

### Basic FSK Encoding/Decoding

```go
package main

import (
    "fmt"
    "github.com/gleicon/go-fsk/fsk"
)

func main() {
    // Create FSK modem with default config
    config := fsk.DefaultConfig()
    modem := fsk.New(config)
    
    // Encode message
    message := "Hello FSK!"
    signal := modem.Encode([]byte(message))
    
    // Decode signal
    decoded := modem.Decode(signal)
    fmt.Printf("Original: %s\n", message)
    fmt.Printf("Decoded:  %s\n", string(decoded))
}
```

### Ultrasonic Communication

```go
// Configure for ultrasonic frequencies (22kHz+)
config := fsk.UltrasonicConfig()
config.BaseFreq = 22000    // 22kHz base frequency
config.FreqSpacing = 500   // 500Hz spacing between symbols
config.BaudRate = 100      // 100 symbols/second

modem := fsk.New(config)

// Now encode/decode at ultrasonic frequencies
signal := modem.Encode([]byte("Secret message"))
```

### Real-Time Transmission

```go
// Create real-time transmitter
transmitter, err := fsk.NewRealTimeTransmitter(modem)
if err != nil {
    log.Fatal(err)
}
defer transmitter.Close()

// Transmit message through speakers
err = transmitter.Transmit([]byte("Live transmission!"))
if err != nil {
    log.Fatal(err)
}
```

### Real-Time Reception

```go
// Create receiver with callback
receiver, err := fsk.NewRealTimeReceiver(modem, func(data []byte) {
    fmt.Printf("Received: %s\n", string(data))
})
if err != nil {
    log.Fatal(err)
}
defer receiver.Close()

// Start listening
err = receiver.Start()
if err != nil {
    log.Fatal(err)
}

// Listen for 10 seconds
time.Sleep(10 * time.Second)
receiver.Stop()
```

### Duplex Chat

```go
// Create chat session
chatSession, err := fsk.NewChatSession(modem)
if err != nil {
    log.Fatal(err)
}
defer chatSession.Close()

// Start duplex communication
err = chatSession.Start()
if err != nil {
    log.Fatal(err)
}

// Send messages
chatSession.SendMessage("Hello from Go!")

// Receive messages
go func() {
    for msg := range chatSession.ReceiveMessages() {
        fmt.Printf("Received: %s\n", msg)
    }
}()
```

### WAV File Operations

```go
// Generate WAV file
signal := modem.Encode([]byte("File message"))
err := modem.WriteWAVFile("output.wav", signal)

// Read WAV file
signal, err := modem.ReadWAVFile("input.wav")
if err == nil {
    decoded := modem.Decode(signal)
    fmt.Printf("From file: %s\n", string(decoded))
}
```

## Configuration Options

### Config Struct

```go
type Config struct {
    BaseFreq    float64 // Base frequency in Hz
    FreqSpacing float64 // Frequency spacing between symbols  
    Order       int     // FSK order (2^n symbols)
    BaudRate    float64 // Symbol rate (symbols per second)
    SampleRate  int     // Audio sample rate
}
```

### Predefined Configurations

```go
// Standard audible FSK
config := fsk.DefaultConfig()
// BaseFreq: 1000Hz, FreqSpacing: 200Hz, Order: 2, BaudRate: 100

// Ultrasonic FSK  
config := fsk.UltrasonicConfig()
// BaseFreq: 22000Hz, FreqSpacing: 500Hz, Order: 2, BaudRate: 100
```

### Custom Configuration Examples

```go
// High-speed 4-FSK (4 bits per symbol)
config := fsk.Config{
    BaseFreq:    2000,  // 2kHz base
    FreqSpacing: 300,   // 300Hz between symbols
    Order:       2,     // 4-FSK (2 bits per symbol) 
    BaudRate:    200,   // 200 symbols/sec = 400 bits/sec
    SampleRate:  48000,
}

// Near-ultrasonic steganography
config := fsk.Config{
    BaseFreq:    18000, // Just at edge of hearing
    FreqSpacing: 400,
    Order:       3,     // 8-FSK (3 bits per symbol)
    BaudRate:    50,    // Slower for reliability
    SampleRate:  48000,
}
```

## Frequency Planning

### Audible Range (Human hearing: 20Hz-20kHz)
- **Low frequencies (1-4kHz)**: Good for noisy environments
- **Mid frequencies (4-8kHz)**: Balanced performance
- **High frequencies (8-18kHz)**: Less interference, some hearing loss

### Ultrasonic Range (>20kHz)
- **Near-ultrasonic (18-22kHz)**: Some people can hear
- **True ultrasonic (22kHz+)**: Inaudible to humans
- **High ultrasonic (25kHz+)**: Excellent steganography

### Symbol Spacing Guidelines
- **Minimum spacing**: 2 × BaudRate (Nyquist criterion)
- **Recommended spacing**: 3-5 × BaudRate (noise margin)
- **High-noise environments**: 10+ × BaudRate

## Performance Characteristics

### Typical Performance Metrics
- **Bit rates**: 50-1000 bits/second
- **Frequency range**: 100Hz - 25kHz  
- **Real-time latency**: ~10ms (audio buffer dependent)
- **SNR requirements**: >10dB for reliable decoding

### Optimization Tips
1. **Higher order FSK** increases data rate but requires better SNR
2. **Lower baud rates** improve reliability in noisy conditions
3. **Ultrasonic frequencies** avoid audible interference
4. **Wider frequency spacing** improves noise immunity

## Error Handling

All functions return Go errors for proper error handling:

```go
transmitter, err := fsk.NewRealTimeTransmitter(modem)
if err != nil {
    log.Fatalf("Failed to create transmitter: %v", err)
}

err = transmitter.Transmit(data)
if err != nil {
    log.Fatalf("Transmission failed: %v", err)
}
```

## Thread Safety

- **Modem encoding/decoding**: Thread-safe
- **Real-time audio objects**: Not thread-safe, use from single goroutine
- **Chat sessions**: Thread-safe for SendMessage() and ReceiveMessages()

## Platform Support

The FSK package uses [Malgo](https://github.com/gen2brain/malgo) for cross-platform audio:

- **Windows**: WASAPI, DirectSound, WinMM
- **macOS**: CoreAudio  
- **Linux**: ALSA, PulseAudio, JACK
- **BSD**: OSS, sndio

## Examples Directory

See the main CLI application (`main.go`) for complete usage examples of all features.
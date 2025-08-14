# GO-FSK

High-order FSK (Frequency Shift Keying) implementation supporting both audible and ultrasonic frequencies with real-time audio capabilities.

This repository contains:
- **FSK Library** (`fsk/`): Reusable Go package for FSK modulation/demodulation
- **CLI Tool** (`cmd/fsk-modem/`): Command-line FSK modem application  
- **Examples** (`examples/`): Demonstration programs and use cases

## Features

### Configurable Parameters

- Base frequency and spacing: Works from audible (1000Hz) to ultrasonic (22kHz+)
- FSK order: 2^n symbols (2, 4, 8, 16 symbols for 1, 2, 3, 4 bits per symbol)
- Baud rate: Adjustable symbol rate
- Sample rate: 48kHz for ultrasonic support

### Operation Modes

- WAV file mode: Generate and decode audio files
- Real-time mode: Live microphone/speaker operation
- Chat mode: Full-duplex communication
- Test mode: Local encode/decode verification

## Quick Start

### Prerequisites

- Go 1.19 or later
- Audio hardware supporting 48kHz sample rate
- For ultrasonic frequencies: speakers/microphones capable of >20kHz

### Installation

```bash
# Clone the repository
git clone https://github.com/gleicon/go-fsk
cd go-fsk

# Quick build and test
make build
make demo

# Or build everything
make all
```

### Alternative Installation Methods

```bash
# Using Go directly
go build -o fsk-modem ./cmd/fsk-modem

# Install to GOPATH/bin
make install

# Or with Go
go install github.com/gleicon/go-fsk/cmd/fsk-modem@latest
```

## Build System

GO-FSK includes a comprehensive Makefile for easy building and development:

```bash
# Show all available commands
make help

# Build main fsk-modem binary
make build

# Build all examples
make examples

# Run tests and checks
make test
make check

# Quick demos
make demo              # Basic FSK test
make demo-file         # WAV file generation/decoding
make demo-ultrasonic   # Ultrasonic example

# Development workflow
make dev               # Format, vet, build, test
make clean             # Clean build artifacts
```

## Usage Examples

### Basic Usage (using binary in build/ directory)

```bash
# Build first
make build

# Generate WAV file
./build/fsk-modem -mode tx -msg "Hello World" -output signal.wav

# Decode WAV file
./build/fsk-modem -mode rx -input signal.wav

# Run built-in test
./build/fsk-modem -test
```

### Using Make Demos

```bash
# Quick test without manual building
make demo

# Generate and decode WAV file
make demo-file

# Run ultrasonic example
make demo-ultrasonic
```

### Real-Time Mode

```bash
# Real-time transmit
./build/fsk-modem -mode rtx -msg "Live message"

# Real-time receive (10 second duration)
./build/fsk-modem -mode rrx -duration 10
```

### Chat Mode

```bash
# Terminal UI chat (with built examples)
make example-chat
./build/examples/chat-tui Alice

# Or using main binary
./build/fsk-modem -mode chat

# Ultrasonic chat (inaudible)
./build/fsk-modem -mode chat -freq "22000,500"
```

### Advanced Examples

```bash
# High-order FSK (4-bit symbols)
./build/fsk-modem -mode rtx -msg "Data" -order 4 -freq "2000,300"

# Ultrasonic beacon
./build/fsk-modem -mode rtx -msg "Secret" -freq "25000,200" -baud 50

# File transfer
./build/fsk-modem -mode tx -file document.txt -output data.wav
./build/fsk-modem -mode rx -input data.wav -file received.txt
```

## Library Usage

The FSK functionality is available as a reusable Go package:

```go
import "github.com/gleicon/go-fsk/fsk"

// Basic usage
config := fsk.DefaultConfig()
modem := fsk.New(config)

// Encode message
signal := modem.Encode([]byte("Hello FSK"))

// Decode signal
decoded := modem.Decode(signal)
fmt.Printf("Decoded: %s\n", string(decoded))

// Ultrasonic configuration
config = fsk.UltrasonicConfig()
config.BaseFreq = 22000
modem = fsk.New(config)

// Real-time transmission
transmitter, err := fsk.NewRealTimeTransmitter(modem)
if err == nil {
    transmitter.Transmit([]byte("Live data"))
    transmitter.Close()
}
```

See `fsk/README.md` for complete library documentation.

## Technical Implementation

### FSK Algorithm

- Encoding: Maps binary data to frequency symbols using sine wave generation
- Decoding: Uses correlation detection to identify transmitted frequencies
- Symbol mapping: MSB-first bit packing for consistent encoding/decoding

### Audio Processing

- Real-time audio I/O via Malgo library
- 48kHz sampling rate for ultrasonic capability
- Configurable symbol periods based on baud rate
- Phase-continuous frequency generation

### Platform Support

- Windows: WASAPI, DirectSound, WinMM
- macOS: CoreAudio
- Linux: ALSA, PulseAudio, JACK
- BSD: OSS, sndio

## Parameters

```
-mode string
    Mode: 'tx' for transmit, 'rx' for receive, 'rtx' for real-time transmit, 
    'rrx' for real-time receive, 'chat' for duplex chat

-msg string
    Message to transmit

-file string
    File to transmit or save received data

-input string
    Input WAV file for receive mode (default "input.wav")

-output string
    Output WAV file for transmit mode (default "output.wav")

-freq string
    Base frequency and spacing in Hz (base,spacing) (default "1000,200")

-order int
    FSK order (2^n symbols, typically 2-4) (default 2)

-baud float
    Symbol rate (symbols per second) (default 100)

-duration float
    Receive duration in seconds (real-time rx mode) (default 5)

-test
    Run test mode (encode then decode)
```

## Frequency Ranges

- Audible: 1000-4000 Hz
- Near-ultrasonic: 18000-22000 Hz
- Ultrasonic: 22000+ Hz (inaudible to humans)

## Applications

- Acoustic data transmission
- Ultrasonic communication
- Audio steganography
- Device-to-device pairing
- IoT sensor networks
- Covert communication channels

## Performance

- Bit rates: 50-1000 bits/second
- Frequency range: 100Hz - 25kHz
- Real-time latency: ~10ms
- SNR requirements: >10dB for reliable decoding

## Examples

The `examples/` directory contains demonstrations of various FSK capabilities:

### Available Examples

```bash
# Build all examples
make examples

# Individual example builds
make example-simple          # Basic FSK encode/decode
make example-ultrasonic      # Ultrasonic frequency demos  
make example-chat           # Terminal UI chat application
make example-frequency-test  # Frequency collision testing
```

### Example Descriptions

- **`simple/`**: Basic FSK encoding/decoding with default settings
- **`ultrasonic/`**: Ultrasonic communication examples (22kHz+, inaudible)
- **`chat-tui/`**: Terminal UI chat application with multi-channel support
- **`frequency-test/`**: Frequency collision and mixing demonstrations

### Running Examples

```bash
# Simple FSK demo
./build/examples/simple

# Ultrasonic demonstration
./build/examples/ultrasonic

# Interactive chat (requires terminal)
./build/examples/chat-tui [username]

# Frequency collision tests
./build/examples/frequency-test [scenario]
```

## Development

### Development Workflow

```bash
# Complete development cycle
make dev                # Format, vet, build, test

# Individual steps
make fmt               # Format source code
make vet               # Run go vet
make lint              # Static analysis (requires golangci-lint)
make test              # Run all tests

# Cleaning
make clean             # Remove all build artifacts
make clean-examples    # Remove only example binaries
```

### Testing

```bash
# Run all tests
make test

# Quick functionality tests
make demo
make test-chat         # Interactive chat test
make test-frequencies  # Frequency collision test
```

### Distribution

```bash
# Build for multiple platforms
make dist

# Create release (requires TAG)
make release TAG=v1.0.0
```

## Troubleshooting

### Build Issues

```bash
# Download dependencies
make deps

# Clean and rebuild
make clean
make all
```

### Audio Issues

1. **No audio device**: Verify audio hardware supports 48kHz sample rate
2. **Ultrasonic not working**: Check speaker/microphone frequency response >20kHz  
3. **Real-time issues**: Reduce background noise, increase distance between devices
4. **Permission errors**: Ensure microphone permissions are granted

### Performance

- For better range: Use higher amplitude frequencies  
- For reliability: Increase frequency spacing
- For speed: Use higher FSK order (more bits per symbol)
- For stealth: Use ultrasonic frequencies (22kHz+)
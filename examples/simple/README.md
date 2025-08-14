# Simple FSK Example

Basic FSK encoding and decoding demonstration using the FSK library.

## Purpose

Demonstrates fundamental FSK operations including message encoding, signal generation, decoding, and WAV file creation. This example shows the core functionality without real-time audio components.

## How to Run

```bash
cd examples/simple
go run main.go
```

## What It Does

1. Creates FSK modem with default configuration (1000Hz base, 200Hz spacing, 4-FSK)
2. Encodes a test message into audio samples
3. Decodes the samples back to text
4. Verifies encoding/decoding accuracy
5. Saves the FSK signal as test.wav

## Expected Output

```
FSK Configuration:
  Base Frequency: 1000 Hz
  Frequency Spacing: 200 Hz
  FSK Order: 2 (2^2 = 4 symbols)
  Baud Rate: 100 symbols/sec
  Sample Rate: 48000 Hz

Original message: Hello, FSK World! 🎵
Encoded to 42240 audio samples
Decoded message:  Hello, FSK World! 🎵
✅ Encoding/decoding successful!

Saving signal to test.wav...
Signal saved to test.wav (play with: ffplay test.wav)
```

## Technical Details

### FSK Configuration
- **Frequency range**: 1000-1600 Hz (audible)
- **Symbol encoding**: 4-FSK (2 bits per symbol)
- **Data rate**: 200 bits per second
- **Audio format**: 48kHz 16-bit mono WAV

### Signal Processing
- **Phase-continuous generation**: Prevents audio clicks
- **MSB-first bit packing**: Standard bit ordering
- **Correlation decoding**: Frequency detection using cross-correlation

## Files Generated

- `test.wav`: FSK-modulated audio file containing the encoded message

## Use Cases

- Understanding FSK fundamentals
- Testing codec accuracy
- Generating test signals
- Educational demonstrations
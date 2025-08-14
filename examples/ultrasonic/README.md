# Ultrasonic FSK Example

Demonstration of ultrasonic FSK communication using frequencies above human hearing range.

## Purpose

Shows FSK operation at ultrasonic frequencies (22kHz+) for covert communication, steganography, and applications where audible signals would be disruptive.

## How to Run

```bash
cd examples/ultrasonic
go run main.go
```

## What It Does

1. Configures FSK modem for 22kHz ultrasonic operation
2. Tests multiple secret messages at ultrasonic frequencies
3. Demonstrates encoding/decoding accuracy at high frequencies
4. Performs real-time ultrasonic transmission test
5. Generates WAV files for each test case

## Expected Output

```
Ultrasonic FSK Modem
====================
Base Frequency: 22000 Hz
Frequencies: 22000, 22500, 23000, 23500 Hz
Human audible: false

Test 1: Secret message 1
  Encoded to 30720 samples (0.64 seconds)
  ✅ Decoded successfully: Secret message 1
  Saved to ultrasonic_test_1.wav

...

Real-time ultrasonic transmission test:
=======================================
Transmitting: Ultrasonic test signal
Frequency range: 22000-23500 Hz (ultrasonic)
✅ Ultrasonic transmission complete!
(Signal was transmitted at 22kHz - inaudible to humans)
```

## Technical Details

### Ultrasonic Configuration
- **Base frequency**: 22000 Hz (above human hearing)
- **Frequency spacing**: 500 Hz (wide spacing for noise immunity)
- **Symbol encoding**: 4-FSK (2 bits per symbol)
- **Data rate**: 200 bits per second
- **Frequency range**: 22000-23500 Hz

### Audio Characteristics
- **Human audibility**: Frequencies above 20kHz are inaudible to most humans
- **Equipment requirements**: Speakers and microphones must support >22kHz
- **Propagation**: Ultrasonic signals have limited range and directional properties
- **Interference**: Less interference from environmental noise

### Real-Time Transmission
- Uses Malgo library for cross-platform audio I/O
- Streams audio directly through system speakers
- Demonstrates live ultrasonic communication capability

## Files Generated

- `ultrasonic_test_1.wav` through `ultrasonic_test_4.wav`: Test signal files
- Each file contains encoded message at ultrasonic frequencies

## Applications

- **Steganographic communication**: Hidden data transmission
- **Device-to-device pairing**: Proximity-based device discovery
- **Covert channels**: Security research and penetration testing
- **Audio watermarking**: Digital rights management
- **IoT sensor networks**: Machine-to-machine communication
- **Accessibility**: Assistive technology for hearing impaired

## Hardware Requirements

- Speakers capable of >22kHz output
- Microphones with >22kHz response
- Audio interface with 48kHz sample rate support

## Testing Ultrasonic Capability

To verify your system supports ultrasonic frequencies:

1. Run the example and check for "Human audible: false"
2. Play generated WAV files - you should hear nothing
3. Use audio spectrum analyzer to confirm 22kHz+ content
4. Test real-time transmission between two devices
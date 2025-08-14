# FSK Core Package

Pure Go implementation of FSK (Frequency Shift Keying) algorithm with no external dependencies.

## Features

- **Zero Dependencies**: Pure Go implementation
- **Platform Independent**: Works on all platforms including WebAssembly
- **Configurable**: Support for different FSK orders, frequencies, and baud rates
- **High Performance**: Optimized signal processing algorithms

## Usage

```go
import "github.com/gleicon/go-fsk/fsk/core"

// Create modem with default configuration
config := core.DefaultConfig()
modem := core.New(config)

// Encode binary data to audio signal
message := []byte("Hello, FSK!")
signal := modem.Encode(message)

// Decode audio signal back to binary data
decoded := modem.Decode(signal)
fmt.Printf("Decoded: %s\n", string(decoded))
```

## Configuration

### Default Configuration
```go
config := core.DefaultConfig()
// BaseFreq: 1000 Hz
// FreqSpacing: 200 Hz  
// Order: 2 (4 symbols)
// BaudRate: 100 symbols/sec
// SampleRate: 48000 Hz
```

### Ultrasonic Configuration
```go
config := core.UltrasonicConfig()
// BaseFreq: 22000 Hz (inaudible to humans)
// FreqSpacing: 500 Hz
// Order: 2 (4 symbols)
// BaudRate: 100 symbols/sec
// SampleRate: 48000 Hz
```

### Custom Configuration
```go
config := core.Config{
    BaseFreq:    2000,    // Base frequency in Hz
    FreqSpacing: 300,     // Frequency spacing between symbols
    Order:       3,       // FSK order (2^3 = 8 symbols)
    BaudRate:    200,     // Symbol rate
    SampleRate:  48000,   // Audio sample rate
}
modem := core.New(config)
```

## API Reference

### Types

#### `Config`
Configuration parameters for FSK modem:
- `BaseFreq float64`: Base frequency in Hz
- `FreqSpacing float64`: Frequency spacing between symbols  
- `Order int`: FSK order (2^n symbols)
- `BaudRate float64`: Symbol rate (symbols per second)
- `SampleRate int`: Audio sample rate

#### `Modem`
FSK modem instance with encoding/decoding capabilities.

### Functions

#### `DefaultConfig() Config`
Returns default FSK configuration suitable for general use.

#### `UltrasonicConfig() Config`  
Returns configuration optimized for ultrasonic communication.

#### `New(config Config) *Modem`
Creates new FSK modem with given configuration.

### Methods

#### `(m *Modem) Config() Config`
Returns the modem's configuration.

#### `(m *Modem) Frequencies() []float64`
Returns array of frequencies used by this modem.

#### `(m *Modem) SymbolPeriod() int`
Returns number of samples per symbol.

#### `(m *Modem) Encode(data []byte) []float32`
Converts binary data to FSK-modulated audio signal.

#### `(m *Modem) Decode(signal []float32) []byte`
Converts FSK-modulated audio signal back to binary data.

## Algorithm Details

### Encoding Process
1. **Symbol Mapping**: Binary data is grouped into symbols based on FSK order
2. **Frequency Assignment**: Each symbol maps to a specific frequency
3. **Signal Generation**: Continuous-phase sine waves generated for each symbol
4. **Output**: Float32 array representing audio samples

### Decoding Process  
1. **Symbol Extraction**: Input signal divided into symbol periods
2. **Correlation Analysis**: Each period tested against all possible frequencies
3. **Symbol Detection**: Frequency with highest correlation selected
4. **Binary Reconstruction**: Symbols converted back to binary data

### Key Features
- **Continuous Phase**: Smooth transitions between frequencies
- **Correlation Detection**: Robust frequency identification
- **MSB-First Encoding**: Consistent bit ordering
- **Configurable Parameters**: Flexible frequency and timing settings

## Performance

### Typical Specifications
- **Bit Rates**: 50-1000 bits/second
- **Frequency Range**: 100Hz - 25kHz  
- **Symbol Rates**: 10-2400 symbols/second
- **Orders**: 1-4 (2-16 symbols)

### Memory Usage
- **Minimal**: No persistent buffers or caches
- **Stateless**: Each encode/decode operation independent
- **Efficient**: Direct float32 operations

## WebAssembly Support

This package compiles to WebAssembly without modifications:

```bash
GOOS=js GOARCH=wasm go build
```

Works in all modern browsers with WebAssembly support.
# FSK Utils Package

Shared utilities for FSK operations, primarily file I/O functionality.

## Features

- **WAV File I/O**: Read and write 16-bit PCM WAV files
- **Platform Independent**: Works on all Go-supported platforms
- **Core Integration**: Seamless integration with FSK core package

## Usage

### Writing WAV Files

```go
import (
    "github.com/gleicon/go-fsk/fsk/core"
    "github.com/gleicon/go-fsk/fsk/utils"
)

// Create FSK signal
modem := core.New(core.DefaultConfig())
signal := modem.Encode([]byte("Hello, FSK!"))

// Write to WAV file
err := utils.WriteWAVFile("output.wav", signal, modem.Config())
if err != nil {
    log.Fatal(err)
}
```

### Reading WAV Files

```go
// Read WAV file
signal, err := utils.ReadWAVFile("input.wav")
if err != nil {
    log.Fatal(err)
}

// Decode FSK signal
modem := core.New(core.DefaultConfig())
decoded := modem.Decode(signal)
fmt.Printf("Decoded: %s\n", string(decoded))
```

## API Reference

### Functions

#### `WriteWAVFile(filename string, signal []float32, config core.Config) error`
Writes audio samples to a WAV file in 16-bit PCM format.

**Parameters:**
- `filename`: Output file path
- `signal`: Audio samples as float32 values (-1.0 to 1.0)
- `config`: FSK configuration for sample rate

**Returns:** Error if file creation or writing fails

#### `ReadWAVFile(filename string) ([]float32, error)`
Reads audio samples from a WAV file and converts to float32 format.

**Parameters:**
- `filename`: Input file path

**Returns:** 
- Audio samples as float32 array (-1.0 to 1.0)
- Error if file reading or parsing fails

## File Format Details

### WAV File Structure
- **Format**: PCM (uncompressed)
- **Bit Depth**: 16-bit signed integers
- **Channels**: 1 (mono)
- **Sample Rate**: Matches FSK configuration
- **Byte Order**: Little-endian

### Conversion Details
- **Float32 to Int16**: `int16(sample * 32767)`
- **Int16 to Float32**: `float32(sample) / 32767.0`
- **Range Clamping**: Values outside [-1.0, 1.0] are clamped
- **Header Size**: Standard 44-byte WAV header

## Error Handling

### Common Errors
- **File Not Found**: Input file doesn't exist
- **Permission Denied**: Insufficient file system permissions
- **Invalid Format**: File is not a valid WAV file
- **Unsupported Format**: Non-16-bit or non-mono WAV files

### Error Examples
```go
signal, err := utils.ReadWAVFile("nonexistent.wav")
if err != nil {
    log.Printf("Failed to read WAV file: %v", err)
}

err = utils.WriteWAVFile("/readonly/output.wav", signal, config)
if err != nil {
    log.Printf("Failed to write WAV file: %v", err)
}
```

## Integration Examples

### CLI Tool Integration
```go
// Generate FSK signal and save to file
signal := modem.Encode(messageData)
err := utils.WriteWAVFile(outputFile, signal, modem.Config())

// Load WAV file and decode FSK signal  
signal, err := utils.ReadWAVFile(inputFile)
if err == nil {
    decoded := modem.Decode(signal)
}
```

### Testing Integration
```go
// Save test signals for analysis
testSignal := modem.Encode([]byte("test message"))
utils.WriteWAVFile("test_output.wav", testSignal, modem.Config())

// Load reference signals for comparison
refSignal, _ := utils.ReadWAVFile("reference.wav")
```

## Performance

### File Size Calculation
- **Formula**: `samples * 2 bytes + 44 byte header`
- **Example**: 10 seconds at 48kHz = 960,044 bytes (~960KB)
- **Compression**: WAV files are uncompressed for maximum compatibility

### Memory Usage
- **Reading**: Entire file loaded into memory as float32 array
- **Writing**: Signal data converted and written in single operation
- **Temporary**: ~4x memory usage during conversion (float32 → int16)

## Platform Compatibility

### Tested Platforms
- **Windows**: All versions with Go support
- **macOS**: All versions with Go support  
- **Linux**: All distributions with Go support
- **BSD**: FreeBSD, OpenBSD, NetBSD

### File System Requirements
- **Read Permissions**: For input WAV files
- **Write Permissions**: For output WAV files and parent directories
- **Disk Space**: Sufficient space for output files
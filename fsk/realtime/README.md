# FSK Realtime Package

Real-time audio I/O functionality for FSK communication using cross-platform audio.

## Features

- **Cross-Platform Audio**: Uses malgo for Windows, macOS, Linux, BSD support
- **Real-Time Processing**: Low-latency audio capture and playback
- **Duplex Communication**: Simultaneous transmit and receive
- **Multi-Channel Support**: Multiple frequency channels
- **Chat Sessions**: Full-duplex communication sessions

## Dependencies

- `github.com/gleicon/go-fsk/fsk/core`: Core FSK algorithm
- `github.com/gen2brain/malgo`: Cross-platform audio I/O

## Usage

### Real-Time Transmission

```go
import (
    "github.com/gleicon/go-fsk/fsk/core"
    "github.com/gleicon/go-fsk/fsk/realtime"
)

modem := core.New(core.DefaultConfig())

// Create transmitter
transmitter, err := realtime.NewTransmitter(modem)
if err != nil {
    log.Fatal(err)
}
defer transmitter.Close()

// Transmit message
err = transmitter.Transmit([]byte("Hello, World!"))
if err != nil {
    log.Fatal(err)
}
```

### Real-Time Reception

```go
// Create receiver with callback
receiver, err := realtime.NewReceiver(modem, func(data []byte) {
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

// Keep listening...
time.Sleep(10 * time.Second)
receiver.Stop()
```

### Full-Duplex Chat

```go
// Create chat session
chatSession, err := realtime.NewChatSession(modem)
if err != nil {
    log.Fatal(err)
}
defer chatSession.Close()

// Start session
err = chatSession.Start()
if err != nil {
    log.Fatal(err)
}

// Send message
chatSession.SendMessage("Hello from chat!")

// Receive messages
go func() {
    for msg := range chatSession.ReceiveMessages() {
        fmt.Printf("Received: %s\n", msg)
    }
}()
```

### Multi-Channel Communication

```go
// Create multi-channel chat
chat := realtime.NewMultiChannelChat("Alice", func(channelID int, username, message string) {
    fmt.Printf("[Channel %d] %s: %s\n", channelID, username, message)
})
defer chat.Close()

// Get predefined channels
channels := realtime.PredefinedChannels()

// Join multiple channels
for _, channel := range channels[:3] {
    err := chat.JoinChannel(channel, 2, 100) // order=2, baud=100
    if err != nil {
        log.Printf("Failed to join channel %d: %v", channel.ID, err)
    }
}

// Send message to specific channel
chat.SendMessage(1, "Hello channel 1!")

// Broadcast to all channels
chat.BroadcastMessage("Hello everyone!")
```

## API Reference

### Types

#### `Transmitter`
Real-time FSK transmitter for audio output.

#### `Receiver`  
Real-time FSK receiver for audio input.

#### `ChatSession`
Full-duplex communication session.

#### `MultiChannelChat`
Multi-channel chat system.

#### `ChannelConfig`
Frequency channel configuration:
- `ID int`: Channel identifier
- `BaseFreq float64`: Base frequency
- `FreqSpacing float64`: Frequency spacing  
- `Name string`: Human-readable name

### Functions

#### `NewTransmitter(modem *core.Modem) (*Transmitter, error)`
Creates new real-time transmitter.

#### `NewReceiver(modem *core.Modem, callback func([]byte)) (*Receiver, error)`
Creates new real-time receiver with message callback.

#### `NewChatSession(modem *core.Modem) (*ChatSession, error)`
Creates new full-duplex chat session.

#### `NewMultiChannelChat(username string, callback func(int, string, string)) *MultiChannelChat`
Creates new multi-channel chat system.

#### `PredefinedChannels() []ChannelConfig`
Returns predefined ultrasonic frequency channels.

#### `DuplexChannels() map[string]struct{TX, RX ChannelConfig}`
Returns duplex channel pairs for point-to-point communication.

### Audio Configuration

The package uses these audio settings:
- **Format**: 16-bit signed integers
- **Channels**: 1 (mono)  
- **Sample Rate**: Matches FSK modem configuration (typically 48kHz)
- **Buffer Size**: Platform-dependent, optimized for low latency

## Platform Support

### Supported Platforms
- **Windows**: WASAPI, DirectSound, WinMM
- **macOS**: CoreAudio
- **Linux**: ALSA, PulseAudio, JACK
- **BSD**: OSS, sndio

### Audio Requirements
- **Sample Rate**: 48kHz recommended for ultrasonic support
- **Bit Depth**: 16-bit minimum
- **Latency**: <10ms for real-time communication
- **Full Duplex**: Simultaneous capture and playback support

## Error Handling

### Common Issues
- **Audio Device Busy**: Close other audio applications
- **Permission Denied**: Grant microphone permissions
- **Sample Rate Mismatch**: Ensure device supports configured sample rate
- **Buffer Underrun**: Reduce system load or increase buffer size

### Troubleshooting
```go
// Check for audio initialization errors
transmitter, err := realtime.NewTransmitter(modem)
if err != nil {
    log.Printf("Audio error: %v", err)
    // Handle gracefully or exit
}
```

## Performance Considerations

### Latency Optimization
- Use lower baud rates for better reliability
- Increase frequency spacing to improve detection
- Minimize system audio buffer sizes
- Close unnecessary audio applications

### CPU Usage
- Higher FSK orders require more processing
- Real-time correlation analysis is CPU-intensive
- Consider lower sample rates for embedded systems

### Memory Usage
- Audio buffers scale with sample rate and buffer duration
- Multi-channel systems use more memory per active channel
- Clean up resources with `Close()` methods

## Thread Safety

All real-time audio operations are thread-safe:
- Internal mutex protection for shared state
- Safe concurrent access to transmitters/receivers  
- Callback functions called from audio threads
- Use channels or mutexes in callback implementations
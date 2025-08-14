# FSK Chat Terminal User Interface

Interactive terminal-based chat application using ultrasonic FSK communication with multi-channel support.

## Purpose

Provides a full-featured chat interface demonstrating practical ultrasonic communication, frequency channel management, and collision avoidance in a user-friendly terminal interface.

## Requirements

- Terminal environment with TTY support
- Audio hardware supporting 48kHz sample rate
- Speakers and microphones with ultrasonic capability (>20kHz)

## How to Run

### Interactive TUI Mode

```bash
cd examples/chat-tui
go run main.go [username]
```

Replace `[username]` with your desired username (optional, defaults to "User").

### Demo Mode (No TTY Required)

```bash
go run demo.go
```

Shows interface mockup and feature overview without requiring terminal capabilities.

## Interface Modes

### Channel Selection Mode (Default)

Navigate and manage frequency channels:

- **Up/Down arrows** or **j/k**: Navigate channel list
- **Enter/Space**: Join selected channel
- **l**: Leave selected channel  
- **c**: Switch to chat mode (requires joined channels)
- **h** or **?**: Show help screen
- **q**: Quit application

### Chat Mode

Send and receive messages:

- **Type message**: Enter text
- **Enter**: Send message to all joined channels
- **Esc**: Return to channel selection
- **Ctrl+C**: Quit application

### Help Mode

View documentation and controls reference.

## Channel Configuration

### Predefined Channels

1. **Channel 1**: 22kHz (ultrasonic, inaudible)
2. **Channel 2**: 24kHz (ultrasonic, inaudible)
3. **Channel 3**: 26kHz (ultrasonic, inaudible)
4. **Channel 4**: 18kHz (near-ultrasonic, some may hear)
5. **Channel 5**: 20kHz (threshold of human hearing)

### Point-to-Point Duplex Pairs

- **Agent A**: TX=22kHz, RX=24kHz
- **Agent B**: TX=24kHz, RX=22kHz

## Technical Implementation

### Multi-Channel Architecture

- **Frequency Division Multiple Access (FDMA)**: Each channel uses separate frequency band
- **Simultaneous channels**: Join multiple channels for broadcast/monitoring
- **Collision avoidance**: 2-3kHz frequency separation prevents interference
- **Real-time audio**: [Malgo](https://github.com/gen2brain/malgo) library provides cross-platform audio I/O

### FSK Parameters

- **Modulation**: 4-FSK (2 bits per symbol)
- **Baud rate**: 100 symbols/second (200 bits/second)
- **Frequency spacing**: 400-500 Hz between symbols
- **Sample rate**: 48kHz (supports ultrasonic frequencies)

### Communication Models

**Broadcast Mode**:

- All users on same channel hear each other
- Similar to IRC/Discord channel
- Collision detection recommended for busy channels

**Point-to-Point Mode**:

- Dedicated TX/RX frequency pairs
- Private communication between two parties
- Full-duplex capability

**Multi-Channel Mode**:

- Monitor multiple channels simultaneously
- Broadcast messages to all joined channels
- Channel-hopping for security

## Example TUI Screens

### Channel Selection

```
┌─────────────────────────────────────────────┐
│             FSK Ultrasonic Chat             │
└─────────────────────────────────────────────┘

Select Frequency Channels:

> 1. Channel 1 (22kHz) [JOINED]
  2. Channel 2 (24kHz)
  3. Channel 3 (26kHz) [JOINED]
  4. Channel 4 (18kHz)
  5. Channel 5 (20kHz)

↑/↓: navigate • enter: join • l: leave • c: chat • q: quit
```

### Chat Interface

```
┌─────────────────────────────────────────────┐
│              Chat Mode - TestUser           │
└─────────────────────────────────────────────┘

Active Channels: Ch1(22kHz), Ch3(26kHz)

[14:30:45] You: Hello ultrasonic world!
[14:30:52] Remote (Ch1): Message received
[14:31:05] OtherUser (Ch3): Hey there!

┌─────────────────────────────────────────────┐
│ Message: Testing ultrasonic FSK chat_       │
└─────────────────────────────────────────────┘

enter: send • esc: back to channels • ctrl+c: quit
```

## Testing Multi-User Communication

### Single Machine Testing

1. Run multiple instances in separate terminals
2. Use different usernames for identification
3. Join same channels for communication
4. Test frequency separation by using different channels

### Multi-Machine Testing

1. Ensure audio hardware supports ultrasonic frequencies
2. Position devices within acoustic range
3. Run chat application on each device
4. Join same frequency channels
5. Verify ultrasonic communication (should be inaudible)

## Troubleshooting

### No Audio I/O

- Verify audio hardware supports 48kHz sample rate
- Check system audio permissions
- Test with lower frequencies first (18-20kHz)

### TTY Errors

- Run in proper terminal environment
- Use demo.go for interface preview
- Ensure terminal supports cursor positioning

### Communication Issues

- Verify ultrasonic capability of speakers/microphones
- Reduce background noise
- Test with audible frequencies first (modify channel config)
- Check frequency separation (minimum 2kHz recommended)

## Applications

**Covert Communication**:

- Security research and penetration testing
- Ultrasonic signals are inaudible to humans
- Difficult to detect without spectrum analysis

**Device Pairing**:

- Proximity-based device discovery
- Audio-based authentication
- IoT device configuration

**Accessibility**:

- Assistive technology integration
- Silent alarm systems
- Machine-to-machine communication

**Educational**:

- Signal processing demonstrations
- Communication protocol teaching
- Frequency planning exercises
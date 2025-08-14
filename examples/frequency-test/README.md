# Frequency Collision and Mixing Test

Comprehensive testing of frequency interference scenarios in FSK communication systems.

## Purpose

Demonstrates how different frequency configurations affect signal quality and decoding accuracy. Tests collision scenarios, interference patterns, and optimal frequency separation strategies.

## How to Run

```bash
cd examples/frequency-test
go run main.go <scenario>
```

### Available Scenarios

```bash
go run main.go 1  # Same frequency collision
go run main.go 2  # Overlapping frequencies  
go run main.go 3  # Separate frequencies (clean)
go run main.go 4  # Multi-channel broadcast
go run main.go 5  # Point-to-point duplex
```

## Scenario Descriptions

### Scenario 1: Same Frequency Collision
Tests two agents transmitting on identical frequencies (22000-23500 Hz).

**Expected Result**: Signal corruption due to destructive interference
```
Agent A: Message from Agent A
Agent B: Message from Agent B
Both using frequencies: 22000-23500 Hz

After collision:
Modem 1 decoded: "Message from Agent B"
Modem 2 decoded: "Message from Agent B"

Results:
Agent A recovery: false
Agent B recovery: true
COLLISION: Both messages corrupted
```

### Scenario 2: Overlapping Frequencies
Tests partial frequency band overlap between two agents.

**Configuration**: 
- Agent A: 22000-23000 Hz
- Agent B: 22500-23500 Hz (500Hz overlap)

### Scenario 3: Separate Frequencies
Tests well-separated frequency bands for clean communication.

**Configuration**:
- Agent A: 22000-23000 Hz  
- Agent B: 24000-25000 Hz (1000Hz separation)

### Scenario 4: Multi-Channel Broadcast
Tests multiple channels with different users broadcasting simultaneously.

**Configuration**:
- Channel 1: 22000 Hz
- Channel 2: 24000 Hz
- Channel 3: 26000 Hz

### Scenario 5: Point-to-Point Duplex
Tests dedicated TX/RX frequency pairs for full-duplex communication.

**Configuration**:
- Agent A: TX=22kHz, RX=24kHz
- Agent B: TX=24kHz, RX=22kHz

## Technical Analysis

### Frequency Interference Types

**Complete Collision (Same Frequencies)**
- Signals add constructively/destructively
- Random data corruption
- One signal may dominate

**Partial Overlap (Close Frequencies)**
- Intermodulation products
- Reduced signal-to-noise ratio
- Partial data recovery possible

**Clean Separation (Distant Frequencies)**
- Minimal cross-talk
- Independent decoding
- Requires sufficient frequency spacing

### Optimal Frequency Planning

**Minimum Separation**: 2 × BaudRate × FrequencySpacing
- For 100 baud, 200Hz spacing: 40kHz minimum
- Practical: 3-5× for noise margin

**Channel Allocation Strategies**:
1. **FDMA**: Fixed frequency division
2. **CSMA**: Collision detection with retry
3. **TDMA**: Time-division multiplexing

## Files Generated

Each scenario generates WAV files for analysis:
- `collision_mixed.wav`: Same frequency interference
- `overlap_mixed.wav`: Partial frequency overlap
- `clean_mixed.wav`: Well-separated channels
- `multichannel_broadcast.wav`: Multiple simultaneous channels
- `duplex_mixed.wav`: Full-duplex communication

## Signal Analysis

Use audio spectrum analyzer tools to examine:
- Frequency content of mixed signals
- Interference patterns
- Signal separation quality

### Recommended Tools
- Audacity (free spectrum analysis)
- GNU Radio (advanced signal processing)
- MATLAB/Octave (mathematical analysis)

## Applications

**Communication System Design**:
- Frequency planning for multi-user systems
- Interference mitigation strategies
- Channel capacity optimization

**Protocol Development**:
- Collision detection algorithms
- Adaptive frequency selection
- Error correction requirements

**Research Areas**:
- Cognitive radio
- Dynamic spectrum access
- Interference characterization
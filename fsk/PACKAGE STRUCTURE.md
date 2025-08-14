# Package Structure

The FSK package has been restructured to separate concerns and enable WebAssembly support.

## New Package Structure

```
fsk/
├── core/           # Pure FSK algorithm (no dependencies)
├── realtime/       # Real-time audio I/O (malgo-based)  
└── utils/          # Shared utilities (WAV file I/O)
```

## Packages and usage

### Packages imports

```go
import (
    "github.com/gleicon/go-fsk/fsk/core"     // For algorithm
    "github.com/gleicon/go-fsk/fsk/realtime" // For audio I/O
    "github.com/gleicon/go-fsk/fsk/utils"    // For WAV files
)
```

### API

**Core Algorithm (Pure FSK):**

```go
config := core.DefaultConfig()
modem := core.New(config)
```

**Real-time Audio:**

```go
transmitter, err := realtime.NewTransmitter(modem)
receiver, err := realtime.NewReceiver(modem, callback)
chatSession, err := realtime.NewChatSession(modem)
```

**WAV File Operations:**

```go
err := utils.WriteWAVFile("output.wav", signal, modem.Config())
signal, err := utils.ReadWAVFile("input.wav")
```

**Channel Management:**
```go
channels := realtime.PredefinedChannels()
chat := realtime.NewMultiChannelChat(username, callback)
```

## Benefits

1. **WebAssembly Support**: Core algorithm works in browsers without malgo dependencies
2. **Cleaner Architecture**: Separation between signal processing and audio I/O
3. **Better Testing**: Core algorithm can be tested independently
4. **Platform Flexibility**: Easier to add other audio backends
5. **Code Reuse**: WASM and CLI share the same core algorithm

package realtime

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/gen2brain/malgo"
	"github.com/gleicon/go-fsk/fsk/core"
)

// ChatSession represents a duplex communication session.
type ChatSession struct {
	modem           *core.Modem
	ctx             *malgo.AllocatedContext
	captureDevice   *malgo.Device
	playbackDevice  *malgo.Device
	capturedSamples []float32
	playbackSignal  []float32
	playbackIndex   uint32
	mu              sync.Mutex
	messageQueue    chan string
	running         bool
}

// NewChatSession creates a new duplex chat session.
func NewChatSession(modem *core.Modem) (*ChatSession, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		// Audio system messages (optional logging)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize audio context: %v", err)
	}

	return &ChatSession{
		modem:        modem,
		ctx:          ctx,
		messageQueue: make(chan string, 10),
	}, nil
}

// Start begins the chat session.
func (c *ChatSession) Start() error {
	// Configure capture device
	captureConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	captureConfig.Capture.Format = malgo.FormatS16
	captureConfig.Capture.Channels = 1
	captureConfig.SampleRate = uint32(c.modem.Config().SampleRate)
	captureConfig.Alsa.NoMMap = 1

	// Configure playback device
	playbackConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	playbackConfig.Playback.Format = malgo.FormatS16
	playbackConfig.Playback.Channels = 1
	playbackConfig.SampleRate = uint32(c.modem.Config().SampleRate)
	playbackConfig.Alsa.NoMMap = 1

	// Capture callback
	onRecvFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		c.mu.Lock()
		defer c.mu.Unlock()

		for i := 0; i < len(pInputSamples); i += 2 {
			sample := int16(binary.LittleEndian.Uint16(pInputSamples[i : i+2]))
			floatSample := float32(sample) / 32767.0
			c.capturedSamples = append(c.capturedSamples, floatSample)
		}

		if len(c.capturedSamples) >= c.modem.SymbolPeriod()*4 {
			decoded := c.modem.Decode(c.capturedSamples)
			if len(decoded) > 0 {
				decodedStr := string(decoded)
				if decodedStr != "" {
					select {
					case c.messageQueue <- decodedStr:
					default:
					}
				}
			}
			c.capturedSamples = c.capturedSamples[:0]
		}
	}

	// Playback callback
	onSendFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		c.mu.Lock()
		defer c.mu.Unlock()

		for i := uint32(0); i < framecount; i++ {
			var sample int16
			if c.playbackIndex < uint32(len(c.playbackSignal)) {
				floatSample := c.playbackSignal[c.playbackIndex]
				if floatSample > 1.0 {
					floatSample = 1.0
				}
				if floatSample < -1.0 {
					floatSample = -1.0
				}
				sample = int16(floatSample * 32767)
				c.playbackIndex++
			} else {
				sample = 0
			}

			outputIndex := i * 2
			if outputIndex+1 < uint32(len(pOutputSample)) {
				binary.LittleEndian.PutUint16(pOutputSample[outputIndex:], uint16(sample))
			}
		}
	}

	// Initialize devices
	captureDevice, err := malgo.InitDevice(c.ctx.Context, captureConfig, malgo.DeviceCallbacks{
		Data: onRecvFrames,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize capture device: %v", err)
	}
	c.captureDevice = captureDevice

	playbackDevice, err := malgo.InitDevice(c.ctx.Context, playbackConfig, malgo.DeviceCallbacks{
		Data: onSendFrames,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize playback device: %v", err)
	}
	c.playbackDevice = playbackDevice

	// Start devices
	err = c.captureDevice.Start()
	if err != nil {
		return fmt.Errorf("failed to start capture device: %v", err)
	}

	err = c.playbackDevice.Start()
	if err != nil {
		return fmt.Errorf("failed to start playback device: %v", err)
	}

	c.running = true
	return nil
}

// SendMessage sends a text message.
func (c *ChatSession) SendMessage(message string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.playbackSignal = c.modem.Encode([]byte(message))
	c.playbackIndex = 0
}

// ReceiveMessages returns a channel for incoming messages.
func (c *ChatSession) ReceiveMessages() <-chan string {
	return c.messageQueue
}

// IsRunning returns true if the session is active.
func (c *ChatSession) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}

// Stop stops the chat session.
func (c *ChatSession) Stop() {
	c.mu.Lock()
	c.running = false
	c.mu.Unlock()

	if c.captureDevice != nil {
		c.captureDevice.Stop()
		c.captureDevice.Uninit()
	}
	if c.playbackDevice != nil {
		c.playbackDevice.Stop()
		c.playbackDevice.Uninit()
	}
}

// Close cleans up resources.
func (c *ChatSession) Close() {
	c.Stop()
	if c.ctx != nil {
		c.ctx.Uninit()
		c.ctx.Free()
	}
	close(c.messageQueue)
}
// Package fsk implements high-order Frequency Shift Keying (FSK) modulation and demodulation.
// It supports both audible and ultrasonic frequencies with configurable parameters.
package fsk

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/gen2brain/malgo"
)

// Config holds the FSK modem configuration parameters.
type Config struct {
	BaseFreq    float64 // Base frequency in Hz
	FreqSpacing float64 // Frequency spacing between symbols
	Order       int     // FSK order (2^n symbols)
	BaudRate    float64 // Symbol rate (symbols per second)
	SampleRate  int     // Audio sample rate
}

// DefaultConfig returns a default FSK configuration.
func DefaultConfig() Config {
	return Config{
		BaseFreq:    1000,
		FreqSpacing: 200,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}
}

// UltrasonicConfig returns a configuration optimized for ultrasonic communication.
func UltrasonicConfig() Config {
	return Config{
		BaseFreq:    22000,
		FreqSpacing: 500,
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}
}

// Modem represents an FSK modem with encoding/decoding capabilities.
type Modem struct {
	config       Config
	symbolPeriod int // Samples per symbol
	frequencies  []float64
	phase        []float64 // Phase accumulators for each frequency
}

// New creates a new FSK modem with the given configuration.
func New(config Config) *Modem {
	modem := &Modem{
		config:       config,
		symbolPeriod: int(float64(config.SampleRate) / config.BaudRate),
		frequencies:  make([]float64, 1<<config.Order), // 2^order frequencies
		phase:        make([]float64, 1<<config.Order),
	}

	// Calculate frequencies for each symbol
	for i := 0; i < len(modem.frequencies); i++ {
		modem.frequencies[i] = config.BaseFreq + float64(i)*config.FreqSpacing
	}

	return modem
}

// Config returns the modem's configuration.
func (m *Modem) Config() Config {
	return m.config
}

// Frequencies returns the array of frequencies used by this modem.
func (m *Modem) Frequencies() []float64 {
	return append([]float64(nil), m.frequencies...) // Return copy
}

// SymbolPeriod returns the number of samples per symbol.
func (m *Modem) SymbolPeriod() int {
	return m.symbolPeriod
}

// Encode converts binary data to FSK-modulated audio signal.
func (m *Modem) Encode(data []byte) []float32 {
	bitsPerSymbol := m.config.Order
	totalBits := len(data) * 8
	symbolCount := (totalBits + bitsPerSymbol - 1) / bitsPerSymbol // Ceiling division

	output := make([]float32, symbolCount*m.symbolPeriod)

	bitIndex := 0
	for symbolIdx := 0; symbolIdx < symbolCount; symbolIdx++ {
		// Extract bits for this symbol
		symbol := 0
		for bit := 0; bit < bitsPerSymbol && bitIndex < totalBits; bit++ {
			byteIdx := bitIndex / 8
			bitInByte := 7 - (bitIndex % 8) // MSB first

			if data[byteIdx]&(1<<bitInByte) != 0 {
				symbol |= 1 << (bitsPerSymbol - 1 - bit)
			}
			bitIndex++
		}

		// Generate waveform for this symbol
		freq := m.frequencies[symbol]
		phaseIncrement := 2 * math.Pi * freq / float64(m.config.SampleRate)

		for sampleIdx := 0; sampleIdx < m.symbolPeriod; sampleIdx++ {
			outputIdx := symbolIdx*m.symbolPeriod + sampleIdx
			if outputIdx < len(output) {
				output[outputIdx] = float32(0.5 * math.Sin(m.phase[symbol]))
				m.phase[symbol] += phaseIncrement

				// Keep phase in range [0, 2π]
				if m.phase[symbol] >= 2*math.Pi {
					m.phase[symbol] -= 2 * math.Pi
				}
			}
		}
	}

	return output
}

// Decode converts FSK-modulated audio signal back to binary data.
func (m *Modem) Decode(signal []float32) []byte {
	symbolCount := len(signal) / m.symbolPeriod
	if symbolCount == 0 {
		return nil
	}

	symbols := make([]int, symbolCount)

	// For each symbol period, determine which frequency has the highest correlation
	for symbolIdx := 0; symbolIdx < symbolCount; symbolIdx++ {
		start := symbolIdx * m.symbolPeriod
		end := start + m.symbolPeriod
		if end > len(signal) {
			end = len(signal)
		}

		maxCorrelation := -1.0
		detectedSymbol := 0

		// Test correlation with each possible frequency
		for freqIdx, freq := range m.frequencies {
			correlation := m.correlateWithFrequency(signal[start:end], freq)
			if correlation > maxCorrelation {
				maxCorrelation = correlation
				detectedSymbol = freqIdx
			}
		}

		symbols[symbolIdx] = detectedSymbol
	}

	// Convert symbols back to bytes
	bitsPerSymbol := m.config.Order
	totalBits := symbolCount * bitsPerSymbol
	byteCount := (totalBits + 7) / 8 // Ceiling division

	output := make([]byte, byteCount)

	bitIndex := 0
	for _, symbol := range symbols {
		for bit := 0; bit < bitsPerSymbol && bitIndex < totalBits; bit++ {
			byteIdx := bitIndex / 8
			bitInByte := 7 - (bitIndex % 8) // MSB first

			if symbol&(1<<(bitsPerSymbol-1-bit)) != 0 {
				output[byteIdx] |= 1 << bitInByte
			}
			bitIndex++
		}
	}

	return output
}

// correlateWithFrequency calculates correlation between signal and reference frequency.
func (m *Modem) correlateWithFrequency(signal []float32, freq float64) float64 {
	phaseIncrement := 2 * math.Pi * freq / float64(m.config.SampleRate)

	var correlation float64
	phase := 0.0

	for _, sample := range signal {
		reference := math.Sin(phase)
		correlation += float64(sample) * reference
		phase += phaseIncrement

		if phase >= 2*math.Pi {
			phase -= 2 * math.Pi
		}
	}

	return math.Abs(correlation) / float64(len(signal))
}

// RealTimeReceiver handles real-time audio capture and decoding.
type RealTimeReceiver struct {
	modem    *Modem
	ctx      *malgo.AllocatedContext
	device   *malgo.Device
	samples  []float32
	mu       sync.Mutex
	callback func([]byte) // Callback for decoded data
}

// NewRealTimeReceiver creates a new real-time receiver.
func NewRealTimeReceiver(modem *Modem, callback func([]byte)) (*RealTimeReceiver, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		// Audio system messages (optional logging)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize audio context: %v", err)
	}

	return &RealTimeReceiver{
		modem:    modem,
		ctx:      ctx,
		callback: callback,
	}, nil
}

// Start begins real-time audio capture and decoding.
func (r *RealTimeReceiver) Start() error {
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = uint32(r.modem.config.SampleRate)
	deviceConfig.Alsa.NoMMap = 1

	onRecvFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		r.mu.Lock()
		defer r.mu.Unlock()

		// Convert int16 samples to float32
		for i := 0; i < len(pInputSamples); i += 2 {
			sample := int16(binary.LittleEndian.Uint16(pInputSamples[i : i+2]))
			floatSample := float32(sample) / 32767.0
			r.samples = append(r.samples, floatSample)
		}

		// Try to decode if we have enough samples
		if len(r.samples) >= r.modem.symbolPeriod*4 { // At least 4 symbols
			decoded := r.modem.Decode(r.samples)
			if len(decoded) > 0 && r.callback != nil {
				r.callback(decoded)
			}
			r.samples = r.samples[:0] // Clear buffer
		}
	}

	device, err := malgo.InitDevice(r.ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onRecvFrames,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize capture device: %v", err)
	}

	r.device = device
	return r.device.Start()
}

// Stop stops the real-time receiver.
func (r *RealTimeReceiver) Stop() {
	if r.device != nil {
		r.device.Stop()
		r.device.Uninit()
	}
}

// Close cleans up resources.
func (r *RealTimeReceiver) Close() {
	r.Stop()
	if r.ctx != nil {
		r.ctx.Uninit()
		r.ctx.Free()
	}
}

// RealTimeTransmitter handles real-time audio generation and playback.
type RealTimeTransmitter struct {
	modem       *Modem
	ctx         *malgo.AllocatedContext
	device      *malgo.Device
	signal      []float32
	sampleIndex uint32
	mu          sync.Mutex
}

// NewRealTimeTransmitter creates a new real-time transmitter.
func NewRealTimeTransmitter(modem *Modem) (*RealTimeTransmitter, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		// Audio system messages (optional logging)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize audio context: %v", err)
	}

	return &RealTimeTransmitter{
		modem: modem,
		ctx:   ctx,
	}, nil
}

// Transmit encodes and transmits data in real-time.
func (t *RealTimeTransmitter) Transmit(data []byte) error {
	t.mu.Lock()
	t.signal = t.modem.Encode(data)
	t.sampleIndex = 0
	t.mu.Unlock()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = uint32(t.modem.config.SampleRate)
	deviceConfig.Alsa.NoMMap = 1

	onSendFrames := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		t.mu.Lock()
		defer t.mu.Unlock()

		for i := uint32(0); i < framecount; i++ {
			var sample int16
			if t.sampleIndex < uint32(len(t.signal)) {
				floatSample := t.signal[t.sampleIndex]
				if floatSample > 1.0 {
					floatSample = 1.0
				}
				if floatSample < -1.0 {
					floatSample = -1.0
				}
				sample = int16(floatSample * 32767)
				t.sampleIndex++
			} else {
				sample = 0
			}

			outputIndex := i * 2
			if outputIndex+1 < uint32(len(pOutputSample)) {
				binary.LittleEndian.PutUint16(pOutputSample[outputIndex:], uint16(sample))
			}
		}
	}

	device, err := malgo.InitDevice(t.ctx.Context, deviceConfig, malgo.DeviceCallbacks{
		Data: onSendFrames,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize playback device: %v", err)
	}
	defer device.Uninit()

	err = device.Start()
	if err != nil {
		return fmt.Errorf("failed to start playback device: %v", err)
	}
	defer device.Stop()

	transmissionTime := float64(len(t.signal)) / float64(t.modem.config.SampleRate)
	time.Sleep(time.Duration(transmissionTime*1000+500) * time.Millisecond)

	return nil
}

// Close cleans up resources.
func (t *RealTimeTransmitter) Close() {
	if t.ctx != nil {
		t.ctx.Uninit()
		t.ctx.Free()
	}
}

// ChatSession represents a duplex communication session.
type ChatSession struct {
	modem           *Modem
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
func NewChatSession(modem *Modem) (*ChatSession, error) {
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
	captureConfig.SampleRate = uint32(c.modem.config.SampleRate)
	captureConfig.Alsa.NoMMap = 1

	// Configure playback device
	playbackConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	playbackConfig.Playback.Format = malgo.FormatS16
	playbackConfig.Playback.Channels = 1
	playbackConfig.SampleRate = uint32(c.modem.config.SampleRate)
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

		if len(c.capturedSamples) >= c.modem.symbolPeriod*4 {
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

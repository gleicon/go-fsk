package realtime

import (
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	"github.com/gen2brain/malgo"
	"github.com/gleicon/go-fsk/fsk/core"
)

// Transmitter handles real-time audio generation and playback.
type Transmitter struct {
	modem       *core.Modem
	ctx         *malgo.AllocatedContext
	device      *malgo.Device
	signal      []float32
	sampleIndex uint32
	mu          sync.Mutex
}

// NewTransmitter creates a new real-time transmitter.
func NewTransmitter(modem *core.Modem) (*Transmitter, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		// Audio system messages (optional logging)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize audio context: %v", err)
	}

	return &Transmitter{
		modem: modem,
		ctx:   ctx,
	}, nil
}

// Transmit encodes and transmits data in real-time.
func (t *Transmitter) Transmit(data []byte) error {
	t.mu.Lock()
	t.signal = t.modem.Encode(data)
	t.sampleIndex = 0
	t.mu.Unlock()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = uint32(t.modem.Config().SampleRate)
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

	transmissionTime := float64(len(t.signal)) / float64(t.modem.Config().SampleRate)
	time.Sleep(time.Duration(transmissionTime*1000+500) * time.Millisecond)

	return nil
}

// Close cleans up resources.
func (t *Transmitter) Close() {
	if t.ctx != nil {
		t.ctx.Uninit()
		t.ctx.Free()
	}
}
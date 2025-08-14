// Package realtime provides real-time audio I/O functionality for FSK communication.
// This package depends on malgo for cross-platform audio support.
package realtime

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/gen2brain/malgo"
	"github.com/gleicon/go-fsk/fsk/core"
)

// Receiver handles real-time audio capture and decoding.
type Receiver struct {
	modem    *core.Modem
	ctx      *malgo.AllocatedContext
	device   *malgo.Device
	samples  []float32
	mu       sync.Mutex
	callback func([]byte) // Callback for decoded data
}

// NewReceiver creates a new real-time receiver.
func NewReceiver(modem *core.Modem, callback func([]byte)) (*Receiver, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		// Audio system messages (optional logging)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize audio context: %v", err)
	}

	return &Receiver{
		modem:    modem,
		ctx:      ctx,
		callback: callback,
	}, nil
}

// Start begins real-time audio capture and decoding.
func (r *Receiver) Start() error {
	deviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.SampleRate = uint32(r.modem.Config().SampleRate)
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
		if len(r.samples) >= r.modem.SymbolPeriod()*4 { // At least 4 symbols
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
func (r *Receiver) Stop() {
	if r.device != nil {
		r.device.Stop()
		r.device.Uninit()
	}
}

// Close cleans up resources.
func (r *Receiver) Close() {
	r.Stop()
	if r.ctx != nil {
		r.ctx.Uninit()
		r.ctx.Free()
	}
}
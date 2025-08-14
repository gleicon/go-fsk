// WebAssembly wrapper for FSK library
package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/gleicon/go-fsk/fsk/core"
)

// Global modem instance
var modem *core.Modem

// FSKConfig represents the configuration for JavaScript
type FSKConfig struct {
	BaseFreq    float64 `json:"baseFreq"`
	FreqSpacing float64 `json:"freqSpacing"`
	Order       int     `json:"order"`
	BaudRate    float64 `json:"baudRate"`
	SampleRate  int     `json:"sampleRate"`
}

// initFSK initializes the FSK modem with given configuration
func initFSK(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf(map[string]interface{}{
			"success": false,
			"error":   "Missing configuration argument",
		})
	}

	// Parse configuration from JavaScript
	configStr := args[0].String()
	var config FSKConfig
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return js.ValueOf(map[string]interface{}{
			"success": false,
			"error":   "Invalid configuration: " + err.Error(),
		})
	}

	// Create FSK configuration
	fskConfig := core.Config{
		BaseFreq:    config.BaseFreq,
		FreqSpacing: config.FreqSpacing,
		Order:       config.Order,
		BaudRate:    config.BaudRate,
		SampleRate:  config.SampleRate,
	}

	// Initialize modem
	modem = core.New(fskConfig)

	return js.ValueOf(map[string]interface{}{
		"success":      true,
		"baseFreq":     config.BaseFreq,
		"freqSpacing":  config.FreqSpacing,
		"order":        config.Order,
		"baudRate":     config.BaudRate,
		"sampleRate":   config.SampleRate,
	})
}

// encodeMessage encodes a text message to FSK audio signal
func encodeMessage(this js.Value, args []js.Value) interface{} {
	if modem == nil {
		return js.ValueOf(map[string]interface{}{
			"success": false,
			"error":   "FSK modem not initialized. Call initFSK first.",
		})
	}

	if len(args) < 1 {
		return js.ValueOf(map[string]interface{}{
			"success": false,
			"error":   "Missing message argument",
		})
	}

	message := args[0].String()
	if message == "" {
		return js.ValueOf(map[string]interface{}{
			"success": false,
			"error":   "Empty message",
		})
	}

	// Encode message to FSK signal
	signal := modem.Encode([]byte(message))

	// Convert float32 slice to JavaScript array
	jsArray := js.Global().Get("Array").New(len(signal))
	for i, sample := range signal {
		jsArray.SetIndex(i, js.ValueOf(float64(sample)))
	}

	return js.ValueOf(map[string]interface{}{
		"success":    true,
		"signal":     jsArray,
		"samples":    len(signal),
		"duration":   float64(len(signal)) / float64(modem.Config().SampleRate),
		"message":    message,
		"sampleRate": modem.Config().SampleRate,
	})
}

// getDefaultConfig returns default FSK configuration
func getDefaultConfig(this js.Value, args []js.Value) interface{} {
	config := core.DefaultConfig()
	return js.ValueOf(map[string]interface{}{
		"baseFreq":    config.BaseFreq,
		"freqSpacing": config.FreqSpacing,
		"order":       config.Order,
		"baudRate":    config.BaudRate,
		"sampleRate":  config.SampleRate,
	})
}

// getUltrasonicConfig returns ultrasonic FSK configuration
func getUltrasonicConfig(this js.Value, args []js.Value) interface{} {
	config := core.UltrasonicConfig()
	return js.ValueOf(map[string]interface{}{
		"baseFreq":    config.BaseFreq,
		"freqSpacing": config.FreqSpacing,
		"order":       config.Order,
		"baudRate":    config.BaudRate,
		"sampleRate":  config.SampleRate,
	})
}

// getModemInfo returns information about the current modem configuration
func getModemInfo(this js.Value, args []js.Value) interface{} {
	if modem == nil {
		return js.ValueOf(map[string]interface{}{
			"success": false,
			"error":   "FSK modem not initialized",
		})
	}

	config := modem.Config()
	frequencies := modem.Frequencies()

	// Convert frequencies to JavaScript array
	jsFreqs := js.Global().Get("Array").New(len(frequencies))
	for i, freq := range frequencies {
		jsFreqs.SetIndex(i, js.ValueOf(freq))
	}

	return js.ValueOf(map[string]interface{}{
		"success":      true,
		"frequencies":  jsFreqs,
		"symbolPeriod": modem.SymbolPeriod(),
		"symbols":      1 << config.Order,
		"baseFreq":     config.BaseFreq,
		"freqSpacing":  config.FreqSpacing,
		"order":        config.Order,
		"baudRate":     config.BaudRate,
		"sampleRate":   config.SampleRate,
	})
}

// generateTone generates a pure sine wave tone for testing
func generateTone(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return js.ValueOf(map[string]interface{}{
			"success": false,
			"error":   "Missing arguments: frequency, duration, sampleRate",
		})
	}

	frequency := args[0].Float()
	duration := args[1].Float()
	sampleRate := int(args[2].Float())

	samples := int(duration * float64(sampleRate))
	signal := make([]float32, samples)

	// Generate sine wave
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		signal[i] = float32(0.5 * js.Global().Get("Math").Call("sin", 2*js.Global().Get("Math").Get("PI").Float()*frequency*t).Float())
	}

	// Convert to JavaScript array
	jsArray := js.Global().Get("Array").New(len(signal))
	for i, sample := range signal {
		jsArray.SetIndex(i, js.ValueOf(float64(sample)))
	}

	return js.ValueOf(map[string]interface{}{
		"success":    true,
		"signal":     jsArray,
		"samples":    len(signal),
		"duration":   duration,
		"frequency":  frequency,
		"sampleRate": sampleRate,
	})
}

func main() {
	// Keep the program running
	c := make(chan struct{}, 0)

	// Register functions to be called from JavaScript
	js.Global().Set("fskInitialize", js.FuncOf(initFSK))
	js.Global().Set("fskEncode", js.FuncOf(encodeMessage))
	js.Global().Set("fskGetDefaultConfig", js.FuncOf(getDefaultConfig))
	js.Global().Set("fskGetUltrasonicConfig", js.FuncOf(getUltrasonicConfig))
	js.Global().Set("fskGetModemInfo", js.FuncOf(getModemInfo))
	js.Global().Set("fskGenerateTone", js.FuncOf(generateTone))

	// Signal that WASM is ready
	js.Global().Call("postMessage", map[string]interface{}{
		"type": "fsk-wasm-ready",
	})

	println("FSK WASM module loaded and ready!")

	// Wait forever
	<-c
}
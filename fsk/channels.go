package fsk

import (
	"fmt"
	"sync"
	"time"
)

// ChannelConfig defines a frequency channel for communication
type ChannelConfig struct {
	ID          int     // Channel identifier
	BaseFreq    float64 // Base frequency for this channel
	FreqSpacing float64 // Frequency spacing
	Name        string  // Human-readable channel name
}

// PredefinedChannels returns common ultrasonic channels
func PredefinedChannels() []ChannelConfig {
	return []ChannelConfig{
		{ID: 1, BaseFreq: 22000, FreqSpacing: 500, Name: "Channel 1 (22kHz)"},
		{ID: 2, BaseFreq: 24000, FreqSpacing: 500, Name: "Channel 2 (24kHz)"},
		{ID: 3, BaseFreq: 26000, FreqSpacing: 500, Name: "Channel 3 (26kHz)"},
		{ID: 4, BaseFreq: 18000, FreqSpacing: 400, Name: "Channel 4 (18kHz)"},
		{ID: 5, BaseFreq: 20000, FreqSpacing: 400, Name: "Channel 5 (20kHz)"},
	}
}

// DuplexChannels returns frequency pairs for point-to-point communication
func DuplexChannels() map[string]struct {
	TX ChannelConfig
	RX ChannelConfig
} {
	return map[string]struct {
		TX ChannelConfig
		RX ChannelConfig
	}{
		"A": {
			TX: ChannelConfig{ID: 1, BaseFreq: 22000, FreqSpacing: 500, Name: "TX-A"},
			RX: ChannelConfig{ID: 2, BaseFreq: 24000, FreqSpacing: 500, Name: "RX-A"},
		},
		"B": {
			TX: ChannelConfig{ID: 2, BaseFreq: 24000, FreqSpacing: 500, Name: "TX-B"},
			RX: ChannelConfig{ID: 1, BaseFreq: 22000, FreqSpacing: 500, Name: "RX-B"},
		},
	}
}

// MultiChannelChat manages communication across multiple frequency channels
type MultiChannelChat struct {
	channels    map[int]*ChatSession
	activeChans []int
	username    string
	mu          sync.RWMutex
	msgCallback func(channelID int, username, message string)
}

// NewMultiChannelChat creates a new multi-channel chat system
func NewMultiChannelChat(username string, msgCallback func(int, string, string)) *MultiChannelChat {
	return &MultiChannelChat{
		channels:    make(map[int]*ChatSession),
		username:    username,
		msgCallback: msgCallback,
	}
}

// JoinChannel joins a specific frequency channel
func (mc *MultiChannelChat) JoinChannel(channelConfig ChannelConfig, order int, baudRate float64) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	config := Config{
		BaseFreq:    channelConfig.BaseFreq,
		FreqSpacing: channelConfig.FreqSpacing,
		Order:       order,
		BaudRate:    baudRate,
		SampleRate:  48000,
	}

	modem := New(config)
	chatSession, err := NewChatSession(modem)
	if err != nil {
		return fmt.Errorf("failed to create chat session for channel %d: %v", channelConfig.ID, err)
	}

	err = chatSession.Start()
	if err != nil {
		return fmt.Errorf("failed to start chat session for channel %d: %v", channelConfig.ID, err)
	}

	// Set up message forwarding
	go func() {
		for msg := range chatSession.ReceiveMessages() {
			if mc.msgCallback != nil {
				mc.msgCallback(channelConfig.ID, "Remote", msg)
			}
		}
	}()

	mc.channels[channelConfig.ID] = chatSession
	mc.activeChans = append(mc.activeChans, channelConfig.ID)

	return nil
}

// LeaveChannel leaves a specific frequency channel
func (mc *MultiChannelChat) LeaveChannel(channelID int) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	chatSession, exists := mc.channels[channelID]
	if !exists {
		return fmt.Errorf("not connected to channel %d", channelID)
	}

	chatSession.Close()
	delete(mc.channels, channelID)

	// Remove from active channels
	for i, id := range mc.activeChans {
		if id == channelID {
			mc.activeChans = append(mc.activeChans[:i], mc.activeChans[i+1:]...)
			break
		}
	}

	return nil
}

// SendMessage sends a message to a specific channel
func (mc *MultiChannelChat) SendMessage(channelID int, message string) error {
	mc.mu.RLock()
	chatSession, exists := mc.channels[channelID]
	mc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("not connected to channel %d", channelID)
	}

	// Add username prefix
	fullMessage := fmt.Sprintf("%s: %s", mc.username, message)
	chatSession.SendMessage(fullMessage)

	return nil
}

// BroadcastMessage sends a message to all active channels
func (mc *MultiChannelChat) BroadcastMessage(message string) error {
	mc.mu.RLock()
	activeChannels := make([]int, len(mc.activeChans))
	copy(activeChannels, mc.activeChans)
	mc.mu.RUnlock()

	var lastErr error
	for _, channelID := range activeChannels {
		if err := mc.SendMessage(channelID, message); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// GetActiveChannels returns list of active channel IDs
func (mc *MultiChannelChat) GetActiveChannels() []int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	channels := make([]int, len(mc.activeChans))
	copy(channels, mc.activeChans)
	return channels
}

// Close closes all channel connections
func (mc *MultiChannelChat) Close() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	for _, chatSession := range mc.channels {
		chatSession.Close()
	}

	mc.channels = make(map[int]*ChatSession)
	mc.activeChans = nil
}

// ChannelAnalyzer analyzes frequency channel conditions
type ChannelAnalyzer struct {
	modem    *Modem
	recorder *RealTimeReceiver
	mu       sync.RWMutex
	activity map[float64]float64 // frequency -> activity level
	running  bool
}

// NewChannelAnalyzer creates a new channel analyzer
func NewChannelAnalyzer() *ChannelAnalyzer {
	// Use wideband configuration for analysis
	config := Config{
		BaseFreq:    18000, // Wide range for analysis
		FreqSpacing: 100,   // Fine resolution
		Order:       2,
		BaudRate:    100,
		SampleRate:  48000,
	}

	modem := New(config)
	return &ChannelAnalyzer{
		modem:    modem,
		activity: make(map[float64]float64),
	}
}

// StartAnalysis begins monitoring channel activity
func (ca *ChannelAnalyzer) StartAnalysis() error {
	receiver, err := NewRealTimeReceiver(ca.modem, func(data []byte) {
		// Activity detected callback
		ca.mu.Lock()
		defer ca.mu.Unlock()

		// Simple activity detection based on received data
		if len(data) > 0 {
			for freq := 18000.0; freq <= 28000.0; freq += 1000.0 {
				ca.activity[freq] = float64(time.Now().Unix())
			}
		}
	})

	if err != nil {
		return err
	}

	ca.recorder = receiver
	ca.running = true
	return receiver.Start()
}

// GetChannelActivity returns current activity levels for frequency ranges
func (ca *ChannelAnalyzer) GetChannelActivity() map[float64]float64 {
	ca.mu.RLock()
	defer ca.mu.RUnlock()

	activity := make(map[float64]float64)
	now := float64(time.Now().Unix())

	for freq, lastActivity := range ca.activity {
		// Activity fades over time
		age := now - lastActivity
		if age < 10.0 { // Active within last 10 seconds
			activity[freq] = 1.0 - (age / 10.0)
		}
	}

	return activity
}

// Stop stops the channel analysis
func (ca *ChannelAnalyzer) Stop() {
	ca.running = false
	if ca.recorder != nil {
		ca.recorder.Close()
	}
}

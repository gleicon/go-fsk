// FSK WASM Demo JavaScript

class FSKDemo {
    constructor() {
        this.audioContext = null;
        this.currentSource = null;
        this.gainNode = null;
        this.wasmReady = false;
        this.currentConfig = null;
        
        this.initializeAudio();
        this.loadWASM();
    }

    async initializeAudio() {
        try {
            // Create AudioContext (handle user gesture requirement)
            this.audioContext = new (window.AudioContext || window.webkitAudioContext)();
            
            // Create gain node for volume control
            this.gainNode = this.audioContext.createGain();
            this.gainNode.connect(this.audioContext.destination);
            this.gainNode.gain.value = 0.5; // 50% volume
            
            console.log('Audio system initialized');
        } catch (error) {
            console.error('Failed to initialize audio:', error);
            this.updateStatus('encodingStatus', 'Error: Failed to initialize audio system', 'error');
        }
    }

    async loadWASM() {
        try {
            // Load and instantiate the WASM module
            const go = new Go();
            const result = await WebAssembly.instantiateStreaming(fetch('fsk.wasm'), go.importObject);
            go.run(result.instance);
            
            // Wait for WASM to signal it's ready
            this.wasmReady = true;
            console.log('WASM module loaded successfully');
            this.updateStatus('configStatus', 'WASM module loaded. Ready to initialize FSK modem.', 'success');
            
            // Load default configuration
            this.loadDefaultConfig();
            
        } catch (error) {
            console.error('Failed to load WASM:', error);
            this.updateStatus('configStatus', 'Error: Failed to load WASM module: ' + error.message, 'error');
        }
    }

    updateStatus(elementId, message, type = '') {
        const element = document.getElementById(elementId);
        if (element) {
            element.textContent = message;
            element.className = 'status ' + type;
        }
    }

    loadDefaultConfig() {
        if (!this.wasmReady) return;
        
        const config = fskGetDefaultConfig();
        this.applyConfigToUI(config);
    }

    loadUltrasonicConfig() {
        if (!this.wasmReady) return;
        
        const config = fskGetUltrasonicConfig();
        this.applyConfigToUI(config);
    }

    loadMSXConfig() {
        // MSX Kansas City Standard configuration
        const config = {
            baseFreq: 1200,
            freqSpacing: 1200,
            order: 1,
            baudRate: 300,
            sampleRate: 48000
        };
        this.applyConfigToUI(config);
    }

    loadSpectrumConfig() {
        // ZX Spectrum-style configuration
        const config = {
            baseFreq: 1000,
            freqSpacing: 1000,
            order: 1,
            baudRate: 1000,
            sampleRate: 48000
        };
        this.applyConfigToUI(config);
    }

    applyConfigToUI(config) {
        document.getElementById('baseFreq').value = config.baseFreq;
        document.getElementById('freqSpacing').value = config.freqSpacing;
        document.getElementById('order').value = config.order;
        document.getElementById('baudRate').value = config.baudRate;
        document.getElementById('sampleRate').value = config.sampleRate;
        
        // Auto-initialize if WASM is ready
        if (this.wasmReady) {
            this.initializeFSK();
        }
    }

    initializeFSK() {
        if (!this.wasmReady) {
            this.updateStatus('configStatus', 'Error: WASM module not ready', 'error');
            return;
        }

        // Get configuration from UI
        const config = {
            baseFreq: parseFloat(document.getElementById('baseFreq').value),
            freqSpacing: parseFloat(document.getElementById('freqSpacing').value),
            order: parseInt(document.getElementById('order').value),
            baudRate: parseFloat(document.getElementById('baudRate').value),
            sampleRate: parseInt(document.getElementById('sampleRate').value)
        };

        // Initialize FSK modem
        try {
            const result = fskInitialize(JSON.stringify(config));
            
            if (result.success) {
                this.currentConfig = config;
                this.updateStatus('configStatus', 'FSK modem initialized successfully!', 'success');
                
                // Get and display frequency information
                const modemInfo = fskGetModemInfo();
                if (modemInfo.success) {
                    this.displayFrequencies(modemInfo.frequencies);
                    this.updateReceiveCommand();
                }
            } else {
                this.updateStatus('configStatus', 'Error: ' + result.error, 'error');
            }
        } catch (error) {
            this.updateStatus('configStatus', 'Error: ' + error.message, 'error');
        }
    }

    displayFrequencies(frequencies) {
        const freqDiv = document.getElementById('frequencyDisplay');
        const freqSpan = document.getElementById('frequencies');
        
        if (frequencies && frequencies.length > 0) {
            const freqList = [];
            for (let i = 0; i < frequencies.length; i++) {
                freqList.push(`Symbol ${i}: ${frequencies[i].toFixed(1)} Hz`);
            }
            freqSpan.textContent = freqList.join(', ');
            freqDiv.style.display = 'block';
        } else {
            freqDiv.style.display = 'none';
        }
    }

    updateReceiveCommand() {
        if (!this.currentConfig) return;
        
        const baseFreq = this.currentConfig.baseFreq;
        const freqSpacing = this.currentConfig.freqSpacing;
        const order = this.currentConfig.order;
        const baudRate = this.currentConfig.baudRate;
        
        const command = `./build/fsk-modem -mode rrx -freq "${baseFreq},${freqSpacing}" -order ${order} -baud ${baudRate} -duration 10`;
        document.getElementById('receiveCommand').value = command;
    }

    async encodeAndPlay() {
        if (!this.wasmReady || !this.currentConfig) {
            this.updateStatus('encodingStatus', 'Error: FSK modem not initialized', 'error');
            return;
        }

        // Resume AudioContext if suspended (user gesture requirement)
        if (this.audioContext.state === 'suspended') {
            await this.audioContext.resume();
        }

        const message = document.getElementById('messageText').value;
        if (!message.trim()) {
            this.updateStatus('encodingStatus', 'Error: Please enter a message', 'error');
            return;
        }

        try {
            // Encode message
            this.updateStatus('encodingStatus', 'Encoding message...', '');
            const result = fskEncode(message);
            
            if (!result.success) {
                this.updateStatus('encodingStatus', 'Error: ' + result.error, 'error');
                return;
            }

            // Convert signal to audio buffer
            const audioBuffer = await this.createAudioBuffer(result.signal, result.sampleRate);
            
            // Play audio
            this.playAudioBuffer(audioBuffer);
            
            // Update status and visualization
            this.updateStatus('encodingStatus', 
                `Encoded "${message}" (${result.samples} samples, ${result.duration.toFixed(2)}s)`, 'success');
            
            this.visualizeWaveform(result.signal);
            this.updateAudioInfo(result);
            
        } catch (error) {
            this.updateStatus('encodingStatus', 'Error: ' + error.message, 'error');
        }
    }

    async createAudioBuffer(signal, sampleRate) {
        const audioBuffer = this.audioContext.createBuffer(1, signal.length, sampleRate);
        const channelData = audioBuffer.getChannelData(0);
        
        // Copy signal data to audio buffer
        for (let i = 0; i < signal.length; i++) {
            channelData[i] = signal[i];
        }
        
        return audioBuffer;
    }

    playAudioBuffer(audioBuffer) {
        // Stop any currently playing audio
        this.stopAudio();
        
        // Create audio source
        this.currentSource = this.audioContext.createBufferSource();
        this.currentSource.buffer = audioBuffer;
        this.currentSource.connect(this.gainNode);
        
        // Add ended event listener
        this.currentSource.onended = () => {
            this.currentSource = null;
            this.updateStatus('encodingStatus', 'Playback finished', 'success');
        };
        
        // Start playback
        this.currentSource.start();
        this.updateStatus('encodingStatus', 'Playing FSK audio...', 'success');
    }

    stopAudio() {
        if (this.currentSource) {
            try {
                this.currentSource.stop();
            } catch (e) {
                // Source may have already stopped
            }
            this.currentSource = null;
        }
    }

    async generateTestTone() {
        if (!this.audioContext) return;
        
        if (this.audioContext.state === 'suspended') {
            await this.audioContext.resume();
        }
        
        try {
            const result = fskGenerateTone(1200, 2.0, this.audioContext.sampleRate);
            
            if (result.success) {
                const audioBuffer = await this.createAudioBuffer(result.signal, result.sampleRate);
                this.playAudioBuffer(audioBuffer);
                
                this.updateStatus('encodingStatus', 
                    `Playing test tone: ${result.frequency} Hz for ${result.duration}s`, 'success');
                this.visualizeWaveform(result.signal.slice(0, 1000)); // Show first 1000 samples
            }
        } catch (error) {
            this.updateStatus('encodingStatus', 'Error generating test tone: ' + error.message, 'error');
        }
    }

    updateVolume(value) {
        if (this.gainNode) {
            this.gainNode.gain.value = value / 100;
        }
        document.getElementById('volumeValue').textContent = value + '%';
    }

    visualizeWaveform(signal) {
        const canvas = document.getElementById('waveformCanvas');
        const ctx = canvas.getContext('2d');
        
        // Clear canvas
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        
        if (!signal || signal.length === 0) return;
        
        // Sample the signal for visualization (max 2000 points)
        const maxPoints = 2000;
        const step = Math.max(1, Math.floor(signal.length / maxPoints));
        const sampledSignal = [];
        
        for (let i = 0; i < signal.length; i += step) {
            sampledSignal.push(signal[i]);
        }
        
        // Draw waveform
        ctx.strokeStyle = '#007acc';
        ctx.lineWidth = 1;
        ctx.beginPath();
        
        const width = canvas.width;
        const height = canvas.height;
        const midY = height / 2;
        
        for (let i = 0; i < sampledSignal.length; i++) {
            const x = (i / sampledSignal.length) * width;
            const y = midY - (sampledSignal[i] * midY);
            
            if (i === 0) {
                ctx.moveTo(x, y);
            } else {
                ctx.lineTo(x, y);
            }
        }
        
        ctx.stroke();
        
        // Draw center line
        ctx.strokeStyle = '#ccc';
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.moveTo(0, midY);
        ctx.lineTo(width, midY);
        ctx.stroke();
    }

    updateAudioInfo(result) {
        const info = `
Signal Info:
- Samples: ${result.samples}
- Duration: ${result.duration.toFixed(2)} seconds
- Sample Rate: ${result.sampleRate} Hz
- Message: "${result.message}"
        `;
        
        document.getElementById('audioInfo').textContent = info;
    }

    copyCommand() {
        const commandInput = document.getElementById('receiveCommand');
        commandInput.select();
        document.execCommand('copy');
        
        // Show feedback
        const originalValue = commandInput.value;
        commandInput.value = 'Copied to clipboard!';
        setTimeout(() => {
            commandInput.value = originalValue;
        }, 1000);
    }
}

// Global functions for HTML onclick handlers
let fskDemo;

function loadDefaultConfig() {
    fskDemo.loadDefaultConfig();
}

function loadUltrasonicConfig() {
    fskDemo.loadUltrasonicConfig();
}

function loadMSXConfig() {
    fskDemo.loadMSXConfig();
}

function loadSpectrumConfig() {
    fskDemo.loadSpectrumConfig();
}

function initializeFSK() {
    fskDemo.initializeFSK();
}

function encodeAndPlay() {
    fskDemo.encodeAndPlay();
}

function stopAudio() {
    fskDemo.stopAudio();
}

function generateTestTone() {
    fskDemo.generateTestTone();
}

function updateVolume(value) {
    fskDemo.updateVolume(value);
}

function copyCommand() {
    fskDemo.copyCommand();
}

// Initialize demo when page loads
document.addEventListener('DOMContentLoaded', () => {
    fskDemo = new FSKDemo();
});
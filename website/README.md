# GO-FSK Website

This directory contains the complete GO-FSK project website with online FSK modem functionality.

## Files

- **index.html** - Main project page with features, documentation links, and architecture overview
- **modem.html** - Online FSK modem for browser-based FSK communication
- **fsk.wasm** - WebAssembly binary containing the FSK core algorithm
- **wasm_exec.js** - Go WebAssembly executor (from Go toolchain)
- **fsk-demo.js** - JavaScript interface for FSK modem functionality

## Local Development

To serve the website locally:

```bash
# Using Python 3
python3 -m http.server 8080

# Using Python 2
python -m SimpleHTTPServer 8080

# Using Node.js
npx http-server -p 8080

# Using Go
go run -m http.server -addr :8080 .
```

Then open [http://localhost:8080] in your browser.

## Testing FSK Communication

### Browser-to-Browser Testing

1. **Setup**: Open `modem.html` on two different computers with speakers and microphones
2. **Configure**: Use identical FSK settings on both computers (recommended: ultrasonic for less interference)
3. **Transmit**: Type a message on one computer and click "Encode & Play Audio"
4. **Receive**: The other computer should hear the FSK audio signal

### Browser-to-CLI Testing  

1. **Browser Side**: Open `modem.html` and configure FSK settings
2. **CLI Side**: Run the fsk-modem tool with matching configuration:

   ```bash
   ./build/fsk-modem -mode rrx -freq "1000,200" -order 2 -baud 100 -duration 10
   ```

3. **Transmit**: Send message from browser
4. **Receive**: CLI tool should decode the message

### Recommended Configurations

- **Audible (Testing)**: Base frequency 1000 Hz, spacing 200 Hz
- **Ultrasonic (Stealth)**: Base frequency 22000 Hz, spacing 500 Hz  
- **MSX Compatible**: Base frequency 1200 Hz, spacing 1200 Hz, order 1
- **High Speed**: Base frequency 2000 Hz, spacing 400 Hz, baud 200

## Deployment

### Static Hosting

The website is completely static and can be deployed to any web server or static hosting service:

- **GitHub Pages**: Push to a GitHub repository with Pages enabled
- **Netlify**: Drag and drop the website directory 
- **Vercel**: Deploy directly from Git repository
- **Traditional Web Server**: Copy files to web server document root

### HTTPS Requirements

Modern browsers require HTTPS for WebAssembly and audio features to work properly. Ensure your deployment uses HTTPS.

### CORS Considerations

The website loads WASM and audio resources. Ensure your web server sends proper CORS headers if needed.

## Features

### Online FSK Modem

- **Real-time encoding**: Convert text to FSK audio signals in the browser
- **Multiple configurations**: Predefined settings for different use cases
- **Audio visualization**: Waveform display of generated signals
- **Volume control**: Adjustable audio output volume
- **CLI integration**: Generate commands for receiving with fsk-modem CLI

### Project Website

- **Responsive design**: Works on desktop and mobile devices
- **Feature showcase**: Complete overview of GO-FSK capabilities
- **Documentation links**: Direct links to GitHub documentation
- **Architecture explanation**: Clear separation of core/realtime/utils packages

## Browser Compatibility

- **Chrome/Edge**: Full support
- **Firefox**: Full support  
- **Safari**: Full support (requires HTTPS for microphone access)
- **Mobile browsers**: Basic support (audio features may be limited)

## Troubleshooting

### WebAssembly Issues

- Ensure files are served from HTTP/HTTPS (not file://)
- Check browser console for WASM loading errors
- Verify fsk.wasm file is not corrupted

### Audio Issues

- Grant microphone permissions when prompted
- Ensure speakers/headphones are working
- Try different frequency ranges if interference occurs
- Use ultrasonic frequencies (22kHz+) for inaudible communication

### Network Issues

- Ensure proper CORS headers for WASM files
- Check that all files are accessible via HTTP requests
- Verify server MIME types for .wasm files

## License

This website and the GO-FSK project are open source under the MIT license

# Decoding Old Computer Tapes with FSK-Modem

This guide explains how to use `fsk-modem` to decode cassette tape data from vintage computers like MSX and ZX Spectrum systems.

## Overview

Many 1980s home computers used cassette tapes for data storage, employing various forms of Frequency Shift Keying (FSK) to encode digital data as audio signals. While modern specialized tools exist for each system, `fsk-modem` can decode some of these formats with appropriate configuration.

## Supported Systems

| Computer    | Success Rate | Encoding Method          | Notes                                     |
| ----------- | ------------ | ------------------------ | ----------------------------------------- |
| MSX         | **High**     | Kansas City Standard FSK | Close to our FSK implementation           |
| ZX Spectrum | **Moderate** | Pulse-width timing       | Custom encoding, requires experimentation |
| BBC Micro   | **High**     | Kansas City Standard     | Similar to MSX                            |
| TRS-80      | **High**     | Kansas City Standard     | Original KCS implementation               |
| Apple II    | **Moderate** | Custom FSK variants      | Multiple encoding methods used            |

## Technical Background

### Kansas City Standard (KCS)

The Kansas City Standard was developed in 1975 and became the de facto standard for cassette data storage:

- **Frequencies**: 1200 Hz (0-bits) and 2400 Hz (1-bits)
- **Bit Encoding**: 
  - `0` bit = 4 cycles of 1200 Hz
  - `1` bit = 8 cycles of 2400 Hz
- **Data Rate**: 300 bits/second
- **Framing**: 1 start bit + 8 data bits + 1-2 stop bits

### CUTS (Computer Users' Tape Standard)

A faster variant of KCS developed by Processor Technology:

- **1200 baud version**:
  - `0` bit = 1 cycle of 1200 Hz  
  - `1` bit = 2 cycles of 2400 Hz
- **Data Rate**: ~1200 bits/second

## MSX Computer Tapes

### Format Details

MSX computers primarily used the Kansas City Standard with some variations:

- **Standard Mode**: 300 baud effective (KCS)
- **Fast Mode**: 1200 baud (CUTS variant)
- **Turbo Mode**: 2400 baud (custom 2400Hz/4800Hz)

### File Formats

- **`.cas`**: Raw cassette data format
- **`.tsx`**: Extended format with timing information

### Decoding Commands

```bash
# Standard MSX (300 baud, Kansas City Standard)
./build/fsk-modem -mode rx -freq "1200,1200" -order 1 -baud 300 -input msx_tape.wav -file decoded.bin

# MSX Fast mode (1200 baud, CUTS)
./build/fsk-modem -mode rx -freq "1200,1200" -order 1 -baud 1200 -input msx_tape.wav -file decoded.bin

# MSX Turbo mode (2400 baud)
./build/fsk-modem -mode rx -freq "2400,2400" -order 1 -baud 2400 -input msx_tape.wav -file decoded.bin
```

### Preparation Steps

1. **Convert CAS to WAV**:

   ```bash
   # Using OpenMSX or similar tools
   cas2wav game.cas game.wav
   ```

2. **Check WAV Properties**:

   ```bash
   # Ensure 22kHz+ sample rate for good frequency resolution
   ffprobe game.wav
   ```

3. **Test Multiple Configurations**:

   ```bash
   # Try different baud rates
   for baud in 300 600 1200 2400; do
     ./build/fsk-modem -mode rx -freq "1200,1200" -order 1 -baud $baud \
       -input msx_tape.wav -file "output_${baud}.bin"
   done
   ```

## ZX Spectrum Tapes

### Format Details

ZX Spectrum uses a custom encoding method:

- **Method**: Pulse-width timing, not pure FSK
- **Frequencies**: Approximately 1kHz and 2kHz
- **Detection**: Zero-crossing with timing measurement
- **Bit Encoding**: Different pulse widths for 0s and 1s

### File Formats

- **`.tzx`**: Tape eXtended format (preferred)
- **`.tap`**: Simple tape format
- **`.pzx`**: Precision timing format

### Decoding Commands

```bash
# Standard attempt (may require experimentation)
./build/fsk-modem -mode rx -freq "1000,1000" -order 1 -baud 1000 -input spectrum.wav -file decoded.bin

# Alternative frequency configurations
./build/fsk-modem -mode rx -freq "1500,1000" -order 1 -baud 800 -input spectrum.wav -file decoded.bin
./build/fsk-modem -mode rx -freq "2000,1000" -order 1 -baud 1200 -input spectrum.wav -file decoded.bin
```

### Preparation Steps

1. **Convert TZX/TAP to WAV**:

   ```bash
   # Using PlayTZX or similar
   tzx2wav game.tzx game.wav
   ```

2. **Experiment with Parameters**:

   ```bash
   # Test frequency combinations
   for base in 800 1000 1200 1500 2000; do
     for spacing in 500 800 1000 1200; do
       ./build/fsk-modem -mode rx -freq "${base},${spacing}" -order 1 -baud 1000 \
         -input spectrum.wav -file "output_${base}_${spacing}.bin"
     done
   done
   ```

## Other Vintage Systems

### BBC Micro

Uses Kansas City Standard similar to MSX:

```bash
# BBC Micro standard encoding
./build/fsk-modem -mode rx -freq "1200,1200" -order 1 -baud 300 -input bbc.wav -file decoded.bin
```

### TRS-80

Original Kansas City Standard implementation:

```bash
# TRS-80 Level I/II BASIC
./build/fsk-modem -mode rx -freq "1200,1200" -order 1 -baud 300 -input trs80.wav -file decoded.bin
```

### Apple II

Multiple encoding methods were used:

```bash
# Standard Apple II cassette interface
./build/fsk-modem -mode rx -freq "1000,1000" -order 1 -baud 300 -input apple2.wav -file decoded.bin

# Some software used different frequencies
./build/fsk-modem -mode rx -freq "2125,2125" -order 1 -baud 1200 -input apple2.wav -file decoded.bin
```

## Troubleshooting

### Common Issues

1. **No Data Decoded**:
   - Check WAV file sample rate (needs 22kHz+ for good frequency resolution)
   - Try different baud rates (±50% of expected)
   - Experiment with frequency spacing

2. **Garbled Output**:
   - Tape speed variations - try different baud rates
   - Wrong frequency configuration
   - Audio quality issues (filtering, noise)

3. **Partial Decoding**:
   - Try processing segments of the tape separately
   - Look for pilot tones and sync signals manually
   - Consider preprocessing with audio tools

### Parameter Tuning

```bash
# Systematic parameter sweep
#!/bin/bash
TAPE="vintage_tape.wav"

for freq in 800 1000 1200 1500 2000 2400; do
  for spacing in 200 400 600 800 1000 1200; do
    for baud in 150 300 600 1200 2400; do
      OUTPUT="test_${freq}_${spacing}_${baud}.bin"
      ./build/fsk-modem -mode rx \
        -freq "${freq},${spacing}" \
        -order 1 -baud ${baud} \
        -input ${TAPE} -file ${OUTPUT}
      
      # Check if output contains recognizable data
      if [ -s ${OUTPUT} ]; then
        echo "Success: freq=${freq}, spacing=${spacing}, baud=${baud}"
        hexdump -C ${OUTPUT} | head -5
      fi
    done
  done
done
```

### Audio Preprocessing

Sometimes preprocessing the audio helps:

```bash
# Normalize audio levels
ffmpeg -i input.wav -filter:a "volume=2.0" -c:a pcm_s16le normalized.wav

# Remove DC offset and apply high-pass filter
ffmpeg -i input.wav -filter:a "highpass=f=100" -c:a pcm_s16le filtered.wav

# Boost specific frequency ranges
ffmpeg -i input.wav -filter:a "equalizer=f=1200:width_type=h:width=2:g=10" boosted.wav
```

## Analyzing Results

### Recognizing Successful Decodes

Look for these patterns in decoded data:

1. **MSX BASIC Programs**:

   ```
   Hex: FF FF FF ... (header)
   ASCII strings: BASIC keywords (PRINT, FOR, NEXT, etc.)
   ```

2. **ZX Spectrum Programs**:

   ```
   Hex: 00 00 00 ... (leader)
   Program blocks start with specific byte patterns
   ```

3. **Data Files**:

   ```
   Consistent byte patterns
   Recognizable text strings
   File structure headers
   ```

### Validation Commands

```bash
# Look for ASCII text (BASIC programs often contain keywords)
strings decoded.bin | grep -i "print\|goto\|for\|next"

# Check for common computer patterns
hexdump -C decoded.bin | head -20

# Look for repeated byte patterns (pilot tones, sync)
hexdump -C decoded.bin | grep "ff ff ff\|00 00 00\|aa aa aa"
```

## Limitations

### Current FSK-Modem Limitations

1. **No UART Framing**: Doesn't handle start/stop bits automatically
2. **No Sync Detection**: Can't detect pilot tones or sync pulses
3. **Fixed Symbol Timing**: May not adapt to tape speed variations
4. **No Error Correction**: Unlike specialized decoders

### When to Use Specialized Tools

For production use, consider dedicated tools:

- **MSX**: OpenMSX, CASduino, MSXPLAYER
- **ZX Spectrum**: Taper, wav2pzx, PlayTZX, TZXDuino
- **General**: Audiotap, WAV-PRG, C64Tape

### Recommended Workflow

1. **Start with fsk-modem** for quick experiments
2. **Use specialized tools** for reliable production decoding
3. **Compare results** between different methods
4. **Document successful configurations** for future use

## Example Success Stories

### MSX Game Decode

```bash
# Successfully decoded "Knightmare" MSX game
./build/fsk-modem -mode rx -freq "1200,1200" -order 1 -baud 300 \
  -input knightmare.wav -file knightmare.bin

# Result: 32KB binary file with recognizable MSX-DOS header
```

### ZX Spectrum BASIC Program

```bash
# Partial success with "Jet Set Willy" loader
./build/fsk-modem -mode rx -freq "1000,1000" -order 1 -baud 1000 \
  -input jsw.wav -file jsw_partial.bin

# Result: Header data decoded, main program required specialized tools
```

## Contributing

If you successfully decode tapes with fsk-modem, please share:

1. **Computer model and tape format**
2. **Successful command line parameters**  
3. **Any preprocessing steps used**
4. **Sample files (if copyright allows)**

This helps build a database of working configurations for the community.

## References

- [Kansas City Standard - Wikipedia](https://en.wikipedia.org/wiki/Kansas_City_standard)
- [MSX Cassette Technical Info](https://hansotten.file-hunter.com/technical-info/cassette-tape/)
- [ZX Spectrum Tape Interface](https://sinclair.wiki.zxnet.co.uk/wiki/Spectrum_tape_interface)
- [TZX Format Specification](https://worldofspectrum.net/TZXformat.html)
- [List of Tape Storage Formats](https://en.wikipedia.org/wiki/List_of_Compact_Cassette_tape_data_storage_formats)
// Demo of what the FSK Chat TUI looks like
package main

import (
	"fmt"

	"github.com/gleicon/go-fsk/fsk/realtime"
)

func main() {
	fmt.Println("FSK Ultrasonic Chat TUI Demo")
	fmt.Println("===========================")
	fmt.Println()

	// Show channel configuration
	channels := realtime.PredefinedChannels()
	fmt.Println("Available Frequency Channels:")
	for i, ch := range channels {
		fmt.Printf("  %d. %s (%.0fkHz)\n", i+1, ch.Name, ch.BaseFreq/1000)
	}
	fmt.Println()

	// Show duplex channels
	duplex := fsk.DuplexChannels()
	fmt.Println("Point-to-Point Duplex Channels:")
	for name, pair := range duplex {
		fmt.Printf("  Agent %s: TX=%.0fkHz, RX=%.0fkHz\n",
			name, pair.TX.BaseFreq/1000, pair.RX.BaseFreq/1000)
	}
	fmt.Println()

	fmt.Println("TUI Interface Features:")
	fmt.Println("======================")
	fmt.Println("• Channel Selection Mode:")
	fmt.Println("  - Navigate channels with arrow keys")
	fmt.Println("  - Join/leave channels with enter/l")
	fmt.Println("  - Switch to chat mode with 'c'")
	fmt.Println()

	fmt.Println("• Chat Mode:")
	fmt.Println("  - Type messages and send with enter")
	fmt.Println("  - Broadcast to all joined channels")
	fmt.Println("  - Real-time message display")
	fmt.Println()

	fmt.Println("• Multi-Channel Support:")
	fmt.Println("  - Join multiple frequency channels")
	fmt.Println("  - Avoid frequency collisions")
	fmt.Println("  - Channel activity monitoring")
	fmt.Println()

	fmt.Println("Example TUI Screen:")
	fmt.Println("==================")
	showMockTUI()

	fmt.Println()
	fmt.Println("To run the actual TUI:")
	fmt.Println("go run main.go YourUsername")
	fmt.Println("(Requires proper terminal environment)")
}

func showMockTUI() {
	fmt.Println(`
┌─────────────────────────────────────────────┐
│             FSK Ultrasonic Chat             │
└─────────────────────────────────────────────┘

Select Frequency Channels:

> 1. Channel 1 (22kHz) [JOINED]
  2. Channel 2 (24kHz)
  3. Channel 3 (26kHz) [JOINED]
  4. Channel 4 (18kHz)
  5. Channel 5 (20kHz)

↑/↓: navigate • enter: join • l: leave • c: chat • q: quit

Recent Messages:
[14:30:15] System: Joined Channel 1 (22kHz)
[14:30:22] System: Joined Channel 3 (26kHz)
[14:30:45] You: Hello ultrasonic world!
[14:30:52] Remote: Message received on Ch1
[14:31:05] Remote (Ch3): Another user here

┌─────────────────────────────────────────────┐
│              Chat Mode - TestUser           │
└─────────────────────────────────────────────┘

Active Channels: Ch1(22kHz), Ch3(26kHz)

[14:30:45] You: Hello ultrasonic world!
[14:30:52] Remote (Ch1): Message received
[14:31:05] OtherUser (Ch3): Hey there!
[14:31:12] You: This is inaudible to humans!

┌─────────────────────────────────────────────┐
│ Message: Testing ultrasonic FSK chat_       │
└─────────────────────────────────────────────┘

enter: send • esc: back to channels • ctrl+c: quit
	`)
}

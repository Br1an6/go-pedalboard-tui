# Go Pedalboard TUI

A terminal-based audio processing interface built with Go, allowing you to apply real-time effects to your audio input.

![alt tag](https://github.com/Br1an6/go-pedalboar-tuid/blob/main/img/tui.png)

## Features

- **Interactive TUI**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for a smooth terminal experience.
- **Device Selection**: Choose your preferred audio input and output devices.
- **Real-time Effects**: Apply various effects including:
  - **Gain**: Simple volume control.
  - **Reverb**: Room simulation.
  - **Distortion**: Waveshaping distortion.
  - **Delay**: Echo effect.
  - **Chorus**: Modulation effect.
  - **Phaser**: Phase shifting.

## Prerequisites

- [Go](https://go.dev/doc/install) 1.25 or later.
- Depending on your system, you may need audio development headers (like PortAudio) for the underlying audio engine.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Br1an6/go-pedalboard-tui.git
   cd go-pedalboard-tui
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

## Usage

1. **Select Input**: Use the arrow keys to navigate the list of available input devices and press **Enter** to select.
2. **Select Output**: Select your desired output device and press **Enter**.
3. **Select Effect**: Choose an effect from the list to start processing.
4. **Control**: The TUI will show the current status. Press `q` or `ctrl+c` to stop the audio stream and exit.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

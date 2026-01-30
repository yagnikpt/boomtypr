# boomtypr

A minimal terminal-based typing test built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- Clean, distraction-free TUI
- Real-time character tracking with visual feedback
- Incorrect characters highlighted with underline
- Responsive text wrapping

## Install

```bash
go install github.com/yagnikpt/boomtypr@latest
```

Or build from source:

```bash
git clone https://github.com/yagnikpt/boomtypr.git
cd boomtypr
make
```

## Usage

```bash
boomtypr
```

## Roadmap

- [x] WPM and accuracy stats display
- [x] Results screen after test completion
- [x] Configurable word count and test duration
- [ ] Multiple word lists
- [x] Test restart functionality

## License

MIT

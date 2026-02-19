# boomtypr

A minimal terminal-based typing test built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).
![Typing Screen](./assets/screen.webp)

## Features

- Clean, distraction-free TUI
- Multiple test modes: Time, Words, and Zen
- Configurable test duration and word count
- Results screen with WPM and accuracy display
- Responsive text wrapping

## Install

### Homebrew (MacOS, Linux)

```bash
brew install yagnikpt/tap/boomtypr
```

### Winget (Windows)

```bash
winget install yagnikbuilds.boomtypr --source winget
```

### Go

```bash
go install github.com/yagnikpt/boomtypr@latest
```

### Build from source

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
- [x] Word erase (ctrl+w / ctlr+backspace)
- [ ] Realtime stats display / better stats

## Inspiration
The UI and flow is inspired by [ashish0kumar/typtea](https://github.com/ashish0kumar/typtea)

## License

MIT

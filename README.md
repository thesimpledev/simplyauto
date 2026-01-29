# SimplyAuto

A Windows automation tool that combines auto-clicking and macro recording functionality.

## Development Approach

This project was developed using AI-assisted programming, but not "vibe coded." The distinction matters: rather than prompting an AI and accepting whatever output it generates, a human developer directed every architectural decision, defined the package structure, specified coding patterns, and reviewed all generated code for correctness and maintainability.

The AI served as a tool to accelerate implementation of well-defined specifications—not as a replacement for software design thinking. Every module boundary, abstraction choice, and structural decision originated from human judgment about what makes code maintainable and extensible.

Multiple LLMs (Claude Opus 4.5 and Gemini) were used throughout development for code review and auditing. Cross-referencing different models helped catch issues that a single model might miss, while human review provided the final judgment on correctness, security, and code quality. This multi-model approach, combined with human oversight, produced more robust code than relying on any single source of review.

## Features

### Auto Clicker
- Configurable click interval (hours, minutes, seconds, milliseconds)
- Optional random offset to vary timing
- Mouse button selection (left, right, middle)
- Single or double click modes
- Repeat a set number of times or run until stopped
- Click at current cursor position or a fixed location

### Macro Recorder
- Record mouse movements, clicks, and keyboard input
- Recordings are kept in memory until you save them
- Playback with adjustable speed (0.5x, 1x, 2x, 4x)
- Loop playback once, a set number of times, or continuously
- Save and load recordings as .simplyauto files

Note: Auto-clicker and macro recorder are mutually exclusive. You can use one or the other, but not both at the same time.

## Hotkeys

| Key | Action |
|-----|--------|
| F6  | Toggle auto clicker |
| F9  | Toggle recording |
| F10 | Toggle playback |
| F11 | Stop |

## Building

The application is built for Windows using GitHub Actions. Clone the repository and push to trigger a build, or run locally on Windows:

```
go build -ldflags="-H windowsgui" -o simplyauto.exe ./cmd/simplyauto
```

## Requirements

- Windows 10 or later
- Go 1.22 or later (for building)

## License

MIT

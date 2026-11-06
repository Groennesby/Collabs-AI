# 𓂀 ANUBIS AGENT CLI (`anubis-agent-cli`)

**Anubis Agent CLI** is the responsive, Hermes-inspired terminal user interface (TUI) and command-line runner written in Go (`go 1.23+`).

## Features
- **Exact Braille Anubis Emblem (`internal/cli/ascii.go`)**: Rendered in Warm Egyptian Gold (`#FFD700`) with glowing Crimson Eyes (`#FF2D2D`).
- **Hermes Classic Color Palette (`internal/cli/theme.go`)**: Deep Midnight Blue-Black (`#1A1A2E`), Warm Gold borders (`#C9A227`), Cornsilk text (`#FFF8DC`), and Lapis Blue (`#1E90FF`), with `/theme` hot-swapping (`anubis-gold`, `lapis-night`, `duat-blood`).
- **Responsive Across All Terminal Sizes (`internal/cli/layout.go`)**:
  - `>= 140 cols`: 2-column desktop layout with full Anubis ASCII art on the left, live telemetry on the right, and Weighing of the Heart + `.synapse` vault sidebar.
  - `100 - 139 cols`: Standard laptop single-column layout.
  - `80 - 99 cols`: Compact split-pane / tmux layout.
  - `< 80 cols` or `< 24 rows`: Ultra-minimal auto-collapsing header with clean word wrapping and zero clipped text.
  - Live `SIGWINCH` terminal resize handler.

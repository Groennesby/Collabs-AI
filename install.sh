#!/usr/bin/env bash
set -e

echo "𓂀 Installing ANUBIS AGENT (annubis & anubis CLI + GUI)..."

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLI_DIR="$SCRIPT_DIR/anubis-agent-cli"
if [ ! -d "$CLI_DIR" ]; then
  CLI_DIR="$SCRIPT_DIR"
fi

if command -v go >/dev/null 2>&1; then
  echo "⚡ Compiling native Go binary (Metal / CUDA accelerated)..."
  (cd "$CLI_DIR" && go build -o bin/annubis . && cp -f bin/annubis bin/anubis)
fi

BIN_SRC="$CLI_DIR/bin/annubis"
if [ ! -f "$BIN_SRC" ]; then
  BIN_SRC="$CLI_DIR/bin/anubis"
fi

TARGET_DIR="/opt/homebrew/bin"
if [ ! -d "$TARGET_DIR" ] || [ ! -w "$TARGET_DIR" ]; then
  if [ -d "/usr/local/bin" ] && [ -w "/usr/local/bin" ]; then
    TARGET_DIR="/usr/local/bin"
  else
    TARGET_DIR="$HOME/.local/bin"
    mkdir -p "$TARGET_DIR"
  fi
fi

ln -sf "$BIN_SRC" "$TARGET_DIR/annubis"
ln -sf "$BIN_SRC" "$TARGET_DIR/anubis"

echo "✓ Installed 'annubis' and 'anubis' to $TARGET_DIR"
echo "𓂀 Type 'annubis' in your terminal to launch Anubis Agent!"

# Strudel Backend Server

A Go WebSocket server for processing Strudel live coding patterns from Neovim.

## Features

- 🎵 Real-time Strudel pattern evaluation
- 🔌 WebSocket communication with Neovim
- 🎛️ Audio event generation and processing
- 📊 Pattern parsing and analysis
- 🔄 Live pattern updates

## Quick Start

### Prerequisites

- Go 1.21 or later
- Neovim with the Strudel-nvim plugin

### Installation

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Build and run the server:
   ```bash
   go run .
   ```

   Or build a binary:
   ```bash
   go build -o strudel-backend
   ./strudel-backend
   ```

3. The server will start on `http://localhost:8080`
   - WebSocket endpoint: `ws://localhost:8080/ws`
   - Health check: `http://localhost:8080/health`

### Usage from Neovim

1. Start the backend server
2. In Neovim, run `:StrudelStart` to connect
3. Write Strudel patterns and evaluate with `:StrudelEvalLine` or `<leader>e`

## Example Patterns

Try these patterns in Neovim:

```javascript
// Basic drum pattern
sound "bd sn hh sn"

// With BPM
sound "bd cp" # bpm 140

// Named pattern
d1 $ sound "bd sn hh cp"

// Musical notes
note "c4 e4 g4 c5"
```

## API

### WebSocket Messages

The server accepts raw Strudel pattern strings and returns JSON responses:

```json
{
  "success": true,
  "pattern": {
    "name": "d1",
    "pattern": "sound \"bd sn hh\"",
    "bpm": 120,
    "elements": ["bd", "sn", "hh"]
  },
  "events": [
    {
      "instrument": "kick",
      "note": "bd",
      "velocity": 0.8,
      "duration": 0.45,
      "time": 0.0
    }
  ],
  "message": "Successfully evaluated pattern 'd1' with 3 events",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Architecture

- `main.go` - WebSocket server and client management
- `strudel.go` - Pattern parsing and audio event generation
- `go.mod` - Go module dependencies

## Next Steps

- [ ] Integrate with SuperCollider for audio synthesis
- [ ] Add MIDI output support
- [ ] Implement pattern scheduling and timing
- [ ] Add more advanced Strudel syntax support
- [ ] Create audio sample management
- [ ] Add pattern recording and playback

## Development

### Testing the Server

```bash
# Test health endpoint
curl http://localhost:8080/health

# Test WebSocket (requires wscat or similar)
wscat -c ws://localhost:8080/ws
```

### Logs

The server provides detailed logging for:
- Client connections/disconnections
- Pattern evaluation
- Audio event generation
- Errors and debugging info


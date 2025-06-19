# Strudel-nvim : IN PRODUCTION

A Neovim plugin for live coding with Strudel patterns, featuring real-time WebSocket communication with a Go backend.

## What is this?

Strudel-nvim lets you evaluate Strudel code directly from Neovim using a high-performance Go WebSocket backend. Write patterns in your editor and hear them instantly!

## Features

🎵 **Live Coding**: Real-time pattern evaluation and audio generation  
🔌 **WebSocket Communication**: Fast, reliable connection between Neovim and Go backend  
🎛️ **Pattern Parsing**: Intelligent parsing of Strudel syntax with audio event generation  
⚡ **Instant Feedback**: Immediate pattern evaluation with detailed logging  
🎹 **Multi-Instrument Support**: Drums, synths, samples, and musical notes  
🎯 **Easy Setup**: Simple installation and configuration  

## Quick Start

### Prerequisites

- Neovim 0.8+
- Go 1.21+
- Lua dependencies will be installed automatically

### Installation

1. **Install the Neovim plugin** (using your preferred plugin manager):

   ```lua
   -- With lazy.nvim
   {
     'yourusername/strudel-nvim',
     config = function()
       require('strudel-nvim').setup({
         -- websocket_url = "ws://localhost:8080/ws", -- Default
         -- keymaps = {                              -- Customize keymaps
         --   start = "<leader>ss",
         --   stop = "<leader>sq", 
         --   eval_line = "<leader>se",
         --   connect = "<leader>sc",
         -- },
         -- disable_keymaps = false,                 -- Set to true to disable all keymaps
       })
     end
   }
   ```

2. **Start the Go backend**:
   ```bash
   ./start-backend.sh
   ```

3. **Connect from Neovim**:
   ```vim
   :StrudelStart
   ```

### Usage

1. Write Strudel patterns in any file
2. Place cursor on a pattern line  
3. Press `<leader>se` or run `:StrudelEval`
4. Listen to your patterns!

## Example Patterns

Try these in Neovim:

```javascript
// Basic drum pattern
sound "bd sn hh sn"

// With BPM control
sound "bd cp" # bpm 140

// Named patterns
d1 $ sound "bd sn hh cp"

// Musical notes
note "c4 e4 g4 c5"

// Rests and variations
sound "bd ~ sn hh"
```

## Commands & Keymaps

| Command | Default Keymap | Description |
|---------|----------------|-------------|
| `:StrudelStart` | `<leader>ss` | Start and connect to the Go backend |
| `:StrudelStop` | `<leader>sq` | Stop and disconnect from backend |
| `:StrudelEval` | `<leader>se` | Evaluate current line |
| `:StrudelConnect` | `<leader>sc` | Connect to existing backend |

### Keymap Customization

You can customize or disable keymaps during setup. See [KEYMAP_CONFIG.md](KEYMAP_CONFIG.md) for detailed configuration options.

## Architecture

```
┌─────────────┐    WebSocket    ┌──────────────┐    Audio
│   Neovim    │◄──────────────► │ Go Backend   │───────►
│  Plugin     │ ws://localhost  │   Server     │
└─────────────┘     :8080/ws    └──────────────┘
```

- **Neovim Plugin** (`lua/`): WebSocket client, commands, and keymaps
- **Go Backend** (`backend/`): WebSocket server, pattern parser, audio engine
- **Communication**: Real-time JSON messages over WebSocket

## Project Status

✅ **Complete**:
- Neovim plugin with WebSocket client
- Go WebSocket server with pattern parsing
- Real-time pattern evaluation
- Audio event generation
- Comprehensive error handling

🚧 **Next Steps**:
- Audio synthesis integration (SuperCollider, Web Audio)
- MIDI output support
- Advanced Strudel syntax
- Pattern scheduling and timing

## Development

### Backend Development

```bash
cd backend
go mod tidy
go run .  # Start development server
```

### Testing

```bash
# Test backend health
curl http://localhost:8080/health

# Test from Neovim
:StrudelStart
:StrudelEvalLine
```

## Contributing

Contributions welcome! Areas of interest:
- Audio engine integration
- Advanced pattern parsing
- Performance optimizations
- Documentation and examples

## License

MIT License - see LICENSE file for details.

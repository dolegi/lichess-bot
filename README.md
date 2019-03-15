# Lichess Bot 
Lichess bot interface for uci compatible engines.

## How to use
- Upgrade to a bot account on lichess. [guide](https://lichess.org/api#operation/botAccountUpgrade)
- Download the binary for your system from the [releases](https://github.com/dolegi/lichess-bot/releases) or compile yourself
- Download a uci compatible engine such as [stockfish](https://stockfishchess.org/download/)
- Create a new toml file for config
- Put bot name, api key and engine path in `config.toml`
- Run with `./lichess-bot config.toml`

## Minimal config
```toml
token = "XXX"
botname = "XXX"
url = "https://lichess.org/api/"

[engine]
path = "path/to/engine"

[challenge]
variants = [
  "standard"
]
speeds = [
  "blitz",
  "bullet"
]
modes = [
  "rated",
  "casual"
]
```

# Config format
```toml
token = "XXX" # API token for lichess
botname = "XXX" # Name of the bot account
url = "https://lichess.org/api/" # Url for lichess api

[engine]
path = "path/to/engine" # Path to uci compatible engine
  [engine.options]
  threads = 1 # Number of CPU threads to use
  hash = 512 # Max memory in MB engine can use
  [engine.go]
  nodes = 1 # Search number of nodes only
  depth = 5 # Search depth limit
  movetime = 5000 # Move time limit in milliseconds

[network]
latency = 100 # Estimated network latency when sending requests in milliseconds

[challenge]
variants = [ # Variants engine supports
  "standard"
]
speeds = [ # Speeds to play at
  "blitz",
  "bullet"
]
modes = [ # Modes to play. Must be "rated" and/or "casual"
  "rated"
]
```

# Releases 
See [releases](https://github.com/dolegi/lichess-bot/releases) for compiled releases
Note: Windows release is untested

# References
- [Lichess docs](https://lichess.org/api#tag/Chess-Bot)
- [python lichess bot](https://github.com/careless25/lichess-bot)
- [uci reference](https://www.shredderchess.com/chess-info/features/uci-universal-chess-interface.html)

# Lichess Bot 
Lichess bot interface for uci compatible engines.

## How to use
- Upgrade to a bot account on lichess. [guide](https://lichess.org/api#operation/botAccountUpgrade)
- `go get github.com/dolegi/lichess-bot` 
- Download a uci compatible engine such as [stockfish](https://stockfishchess.org/download/)
- Create a new toml file for config
- Put bot name, api key and engine path in `config.toml`
- Run with `./lichess-bot config.toml`

### Upgrade to bot account
As well as normal usage shown above, you can use this library to upgrade your account

> WARNING: Upgrading your account cannot be undone

To upgrade run like this
```
./lichess-bot config.toml upgrade
```

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

# Docker
Example usage for using Docker:
`docker build . -t lichess-bot  && docker run -v $(pwd)/config.toml:/app/config.toml lichess-bot ./config.toml`

# Releases 
To install run `go get github/dolegi/lichess-bot`
Note: Windows release is untested

# References
- [Lichess docs](https://lichess.org/api#tag/Chess-Bot)
- [python lichess bot](https://github.com/careless25/lichess-bot)
- [uci reference](https://www.shredderchess.com/chess-info/features/uci-universal-chess-interface.html)

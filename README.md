# shared-shell

basically a reverse shell but the "victim" has access to the shell \
yes, i could have made a terminal multiplexer, or just used one \
but i havent used Go before so why not \
\
also, this has no encryption whatsoever :smiley:

## Building
NOTE: everything for the client must be configured before building
| Command              | Action     |
|----------------------|------------|
| `go build client.go` | For client |
| `go build server.go` | For server |

## Usage
the client is to be run on the machine where the terminal is going to be shared

| Command                      | Action                                  |
|------------------------------|-----------------------------------------|
| `client`                     | Connects to server in client.go         |
| `server [-p port] [-i ip]`   | Hosts server (default: 0.0.0.0:37591)   |
# shared-shell

## This project is incomplete.

basically a reverse shell but the "victim" has access to the shell \
yes, i could have made a terminal multiplexer, or just used one \
but i havent used Go before so why not

## Building
| Command              | Action     |
|----------------------|------------|
| go build `client.go` | For client |
| go build `server.go` | For server |

## Usage
the client is to be run on the machine where the terminal is going to be shared. \
port 37591 must be exposed on the server

## Changing port
No \
..what are you doing with port 37591?
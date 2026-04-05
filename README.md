# tcprcon-cli

- [tcprcon-cli](#tcprcon-cli)
  - [Features](#features)
  - [Installation](#installation)
    - [Binary](#binary)
    - [Docker](#docker)
      - [Basic Usage](#basic-usage)
      - [Persistent Configuration (Profiles)](#persistent-configuration-profiles)
      - [Shell Alias (Optional)](#shell-alias-optional)
    - [Go](#go)
    - [Windows](#windows)
  - [Usage](#usage)
    - [Interactive Mode](#interactive-mode)
    - [Single Command Mode](#single-command-mode)
    - [Keepalive (Pulse)](#keepalive-pulse)
    - [Using Environment Variable for Password](#using-environment-variable-for-password)
  - [Configuration Profiles](#configuration-profiles)
    - [Saving a Profile](#saving-a-profile)
    - [Loading a Profile](#loading-a-profile)
  - [CLI Flags](#cli-flags)
  - [Interactive UX](#interactive-ux)
  - [Protocol Compliance](#protocol-compliance)
  - [Local Development](#local-development)
    - [Prerequisites](#prerequisites)
    - [Getting Started](#getting-started)
  - [Using as a Library](#using-as-a-library)
  - [License](#license)


A fully native RCON client implementation, zero third parties*

<sub>*except for other golang maintained packages about terminal emulators, until i fully master tty :(</sub>

![tcprcon-cli demo](.meta/demo.png)

## Features

- **Interactive Terminal UI**: full-screen exclusive TUI (like vim or nano) with command history and scrollable output
- **Single Command Mode**: execute a single RCON command and exit
- **Multiple Authentication Methods**: supports password via CLI flag, environment variable (`rcon_password`), or secure prompt
- **Keepalive (Pulse)**: configurable periodic command to keep the connection alive on idle servers
- **Configurable Logging**: syslog-style severity levels for debugging
- **Installable as library**: use the RCON client in your own Go projects, ([see examples](#using-as-a-library))

## Installation

### Binary

Linux binaries are available on the [releases page](https://github.com/UltimateForm/tcprcon-cli/releases/latest).

### Docker

The Docker image is pulled automatically on first run, so no separate installation step is needed. Just run the container with your desired flags.

#### Basic Usage

All flags and commands from the main [Usage](#usage) section apply here—just prefix them with `docker run`. For example:

```bash
docker run -it ghcr.io/ultimateform/tcprcon-cli:latest --address=192.168.1.100 --port=7778
```

**Note on Local Servers:** If running the RCON server on the same machine as the Docker container, you need to use either `--network=host` or `host.docker.internal`:

```bash
# Option 1: Use host network
docker run -it --network=host ghcr.io/ultimateform/tcprcon-cli:latest --address=localhost --port=7778
```

```bash
# Option 2: Use host.docker.internal (without --network=host)
docker run -it ghcr.io/ultimateform/tcprcon-cli:latest --address=host.docker.internal --port=7778
```

#### Persistent Configuration (Profiles)

**Note:** `tcprcon-cli` supports configuration profiles out of the box (see [Configuration Profiles](#configuration-profiles)). However, when using Docker, profiles are stored inside the container and lost when it exits. To persist profiles across container runs, use a Docker named volume:

```bash
docker run -it \
  -v tcprcon-config:/root/.config/tcprcon \
  ghcr.io/ultimateform/tcprcon-cli:latest \
  --address=192.168.1.100 --port=7778 --save="my_server"
```

Then load the profile in future runs:

```bash
docker run -it --rm \
  -v tcprcon-config:/root/.config/tcprcon \
  ghcr.io/ultimateform/tcprcon-cli:latest \
  --profile="my_server"
```

#### Shell Alias (Optional)

For convenience, create a shell alias in your `~/.bashrc` or `~/.zshrc`:

```bash
# feel free to name it anything else
alias tcprcon-cli='docker run -it --rm -v tcprcon-config:/root/.config/tcprcon ghcr.io/ultimateform/tcprcon-cli:latest'
```

Then reload your shell and use it:

```bash
source ~/.bashrc  # or ~/.zshrc
tcprcon-cli --address=192.168.1.100 --port=7778 --save="my_server"
tcprcon-cli --profile="my_server"
```

### Go

```bash
go install github.com/UltimateForm/tcprcon-cli@latest
```

Or build from source:

<sub>note: requires golang 1.22+</sub>

```bash
git clone https://github.com/UltimateForm/tcprcon-cli.git
cd tcprcon-cli
go build -o tcprcon-cli .
```

### Windows

Windows is not natively supported due to Unix-specific terminal dependencies. Use [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install) or the Docker image above.

Might change my mind about supporting windows in the future but given that this is essentially a CLI app windows is kind of an after thought for me... and there's likely a better UI app for windows out there.

## Usage

### Interactive Mode

```bash
tcprcon-cli --address=192.168.1.100 --port=7778
```

### Single Command Mode

```bash
tcprcon-cli --address=192.168.1.100 --cmd="playerlist"
```

### Keepalive (Pulse)

To keep the connection alive on idle servers, use `-pulse` with a command your server accepts as a no-op:

```bash
tcprcon-cli --address=192.168.1.100 --pulse="alive"
```

The default interval is 60 seconds. Override it with `-pulse-interval`:

```bash
tcprcon-cli --address=192.168.1.100 --pulse="alive" --pulse-interval=30s
```

Pulse settings can also be saved to a profile:

```bash
tcprcon-cli --address=192.168.1.100 --pulse="alive" --pulse-interval=30s --save="my_server"
```

### Using Environment Variable for Password

```bash
export rcon_password="your_password"
tcprcon-cli --address=192.168.1.100
```

## Configuration Profiles

`tcprcon-cli` supports saving and loading connection profiles to a local configuration file, located at `~/.config/tcprcon/config.json` on Linux/macOS or `%AppData%\tcprcon\config.json` on Windows.

### Saving a Profile

You can save your current connection parameters (address, port, and optionally password) to a named profile using the `--save` flag.

```bash
# Connect to a server and save its details as "my_server"
tcprcon-cli --address=192.168.1.100 --port=7778 --pw="mysecret" --save="my_server"
```
When saving, you will be prompted if you wish to store the password. If you choose 'y', the password will be saved in plaintext within the `config.json` file with restricted `0600` file permissions (read/write only by owner). If you choose 'n' or omit the password, you will be prompted for it when loading the profile.

### Loading a Profile

To load a previously saved profile, use the `--profile` flag.

```bash
# Load the "my_server" profile
tcprcon-cli --profile="my_server"
```

Explicit CLI flags will always override values from a loaded profile. For example:
```bash
# Load "my_server" but connect to a different port
tcprcon-cli --profile="my_server" --port=27015
```

## CLI Flags

```
  -address string
    	RCON address, excluding port (default "localhost")
  -cmd string
    	command to execute, if provided will not enter into interactive mode
  -log uint
    	sets log level (syslog severity tiers) for execution (default 4)
  -port uint
    	RCON port (default 7778)
  -profile string
    	loads a saved profile by name, overriding default flags but overridden by explicit flags.
  -pulse string
    	keepalive method: a command sent on a schedule to keep the connection alive
  -pulse-interval duration
    	keepalive interval, use Go duration format e.g. 30s, 2m (default 1m0s)
  -pw string
    	RCON password, if not provided will attempt to load from env variables, if unavailable will prompt
  -save string
    	saves current connection parameters as a profile. Value is the profile name.
```

## Interactive UX

The interactive terminal UI supports the following keyboard controls:

| Key | Action |
|-----|--------|
| `Enter` | Submit command |
| `Backspace` | Delete last character |
| `↑` / `↓` | Navigate command history |
| `Page Up` | Scroll output up one page |
| `Page Down` | Scroll output down one page |

When scrolled up, a `[↑ N]` indicator is shown in the prompt line, where `N` is the number of lines scrolled above the bottom. Submitting a command snaps the view back to the bottom.

## Protocol Compliance

While `tcprcon-cli` follows the standard Source RCON Protocol, some game servers (like Rust) have non-standard implementations that might introduce unexpected behaviors, such as duplicated responses or incorrect packet IDs, the cli should still work, you might just have to deal with an overly chatty server.

For a detailed breakdown of known server quirks and how they are handled, see the [Caveats section in the core library documentation](https://github.com/UltimateForm/tcprcon#caveats).


## Local Development

You can use the provided `Makefile` and `compose.yaml` to spin up a local development environment. This will start a Mordhau game server in a container.

### Prerequisites

- Docker and Docker Compose
- Go 1.22+
- Make (optional, but recommended)

### Getting Started

1. **Start a Game Server**:

  Beware that the first time you build and run the server it might take a while for its RCON port to be usable, not sure why, but Rust one for example took a few minutes before it was responding, idk.

   ```bash
   make lift-mh-server 
   ```
   
   or 
    
   ```bash
   make lift-rust-server
   ```
   *Note: The server uses `network_mode: host` and may take a few minutes to fully initialize, make sure network_mode is supported by your docker engine*

2. **Build and Run the Client**:
   ```bash
   make run
   ```
   This will build the `tcprcon-cli` binary into `.out/` and execute it against the local server (see step 1) using the default development credentials.

3. **Run Tests**:
   ```bash
   make test
   ```

4. **Dockerized Client**:
   If you prefer to run the client itself inside a container:
   ```bash
   make run-docker
   ```


## Using as a Library

See https://github.com/UltimateForm/tcprcon

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**. See [LICENSE](LICENSE) for details.

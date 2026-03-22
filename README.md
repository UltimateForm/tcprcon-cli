# tcprcon-cli

- [tcprcon-cli](#tcprcon-cli)
  - [Local Development](#local-development)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Interactive Mode](#interactive-mode)
    - [Single Command Mode](#single-command-mode)
    - [Using Environment Variable for Password](#using-environment-variable-for-password)
  - [Configuration Profiles](#configuration-profiles)
  - [CLI Flags](#cli-flags)
  - [Protocol Compliance](#protocol-compliance)
  - [Using as a Library](#using-as-a-library)
    - [Streaming Responses](#streaming-responses)
  - [License](#license)


A fully native RCON client implementation, zero third parties*

<sub>*except for other golang maintained packages about terminal emulators, until i fully master tty :(</sub>

![tcprcon-cli demo](.meta/demo.png)

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

## Features

- **Interactive Terminal UI**: full-screen exclusive TUI (like vim or nano)
- **Single Command Mode**: execute a single RCON command and exit
- **Multiple Authentication Methods**: supports password via CLI flag, environment variable (`rcon_password`), or secure prompt
- **Configurable Logging**: syslog-style severity levels for debugging
- **Installable as library**: use the RCON client in your own Go projects, ([see examples](#using-as-a-library))

## Installation

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

## Usage

### Interactive Mode

```bash
tcprcon-cli --address=192.168.1.100 --port=7778
```

### Single Command Mode

```bash
tcprcon-cli --address=192.168.1.100 --cmd="playerlist"
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
  -pw string
    	RCON password, if not provided will attempt to load from env variables, if unavailable will prompt
  -save string
    	saves current connection parameters as a profile. Value is the profile name.
```

## Protocol Compliance

While `tcprcon-cli` follows the standard Source RCON Protocol, some game servers (like Rust) have non-standard implementations that might introduce unexpected behaviors, such as duplicated responses or incorrect packet IDs, the cli should still work, you might just have to deal with an overly chatty server.

For a detailed breakdown of known server quirks and how they are handled, see the [Caveats section in the core library documentation](https://github.com/UltimateForm/tcprcon#caveats).

## Using as a Library

See https://github.com/UltimateForm/tcprcon

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**. See [LICENSE](LICENSE) for details.

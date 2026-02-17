# tcprcon-cli

- [tcprcon-cli](#tcprcon-cli)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Interactive Mode](#interactive-mode)
    - [Single Command Mode](#single-command-mode)
    - [Using Environment Variable for Password](#using-environment-variable-for-password)
  - [CLI Flags](#cli-flags)
  - [Using as a Library](#using-as-a-library)
    - [Streaming Responses](#streaming-responses)
  - [License](#license)


A fully native RCON client implementation, zero third parties*

<sub>*except for other golang maintained packages about terminal emulators, until i fully master tty :(</sub>

![tcprcon-cli demo](.meta/demo.png)

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
  -pw string
    	RCON password, if not provided will attempt to load from env variables, if unavailable will prompt
```

## Using as a Library

See https://github.com/UltimateForm/tcprcon

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**. See [LICENSE](LICENSE) for details.

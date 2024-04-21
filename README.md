# Game of Life

This project implements Conway's Game of Life as a go cli
## Overview

The system uses BubbleTea and protobufs

## Features

## Getting Started

### Prerequisites

- Go installed on your machine
- `protoc` tooling installed on your machine
  - Windows with Chocolatey: `choco install protoc`

### Installation

1. Clone the repository: `git clone https://github.com/LeahGJoyce/Game-of-Life.git`
2. Navigate into the project directory: `cd Game-of-Life`
3. Compile proto files using `protoc -I="proto/" --go_out="." proto/*.proto`

### Usage

1. run `go run .`

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your suggested changes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

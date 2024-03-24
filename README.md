# Game of Life

This project implements Conway's Game of Life as a microservice, utilizing Go for the service logic and RabbitMQ for messaging. It provides a dynamic web interface to visualize the game's progression in real-time.

## Overview

The system comprises a Dockerized environment with two main components:

- **Go Microservice**: Manages the game logic, starting state, board size, and number of ticks. It publishes each game tick to a RabbitMQ topic.
- **RabbitMQ Server**: Acts as a message broker, receiving game states from the Go microservice and making them available to the web interface.

## Features

- Start a new game with a custom initial state, board size, and tick count.
- Real-time visualization of game states for each tick.
- Utilizes bitboards for efficient state management.

## Getting Started

### Prerequisites

- Docker and Docker Compose installed on your machine.

### Installation

1. Clone the repository: `git clone https://github.com/LeahGJoyce/Game-of-Life.git`
1. Navigate into the project directory: `cd Game-of-Life`
1. Build and start the services with Docker Compose: `docker-compose up --build`

### Usage

- Access the web interface at `http://localhost` to view and start games.
- Use the `/start` endpoint to initiate a new game, specifying the starting state, board size, and tick count.

## Architecture

**These diagrams are a work in progress and may not correctly represent the system architecture.**

## C4 System Context Diagram
```mermaid
C4Context
  title System Context diagram for Game of Life Microservice
  Person(user, "User", "Interacts with the web interface to view and start games.")
  System_Ext(docker, "Docker", "Hosts the containerized microservices environment.")
  System(game_of_life, "Game of Life Microservice", "Implements Conway's Game of Life logic and serves the web interface.")
  
  Rel(user, game_of_life, "Views game via web interface")
  Rel(docker, game_of_life, "Hosts")
```

### C4 Container Diagram
```mermaid
C4Container
  title Container diagram for Game of Life Microservice
  Person(user, "User", "Interacts with the web interface to view and start games.")
  System_Boundary(game_of_life, "Game of Life Microservice") {
    Container(go_microservice, "Go Microservice", "Go", "Manages game logic, state, and communicates with RabbitMQ.")
    Container(rabbitmq_server, "RabbitMQ Server", "Messaging", "Handles messaging for game state updates.")
    ContainerDb(web_interface, "Web Interface", "HTML/CSS/JavaScript", "Provides a dynamic interface to visualize the game.")
    
    Rel(user, web_interface, "Uses")
    Rel(go_microservice, rabbitmq_server, "Publishes game state to")
    Rel(rabbitmq_server, web_interface, "Sends game state updates to", "WebSocket")
  }

```

### C4 Component Diagram
```mermaid
C4Component
  title Component diagram for Go Microservice - Game of Life
  Container(go_microservice, "Go Microservice", "Go", "Manages game logic and communicates with RabbitMQ.")
  Component(game_logic, "Game Logic Component", "Go", "Implements the logic for Conway's Game of Life.")
  Component(message_publisher, "Message Publisher", "Go", "Publishes game state updates to RabbitMQ.")
  Component(web_api, "/start Endpoint", "Go", "Receives requests to start a new game.")
  ContainerDb(rabbitmq_server, "RabbitMQ Server", "Messaging", "Handles messaging for game state updates.")

  Rel(web_api, game_logic, "Initiates game with parameters")
  Rel(game_logic, message_publisher, "Sends game state updates")
  Rel(message_publisher, rabbitmq_server, "Publishes to")
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your suggested changes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

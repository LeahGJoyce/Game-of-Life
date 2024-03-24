package game

import (
	"fmt"

	"github.com/google/uuid"
)
type GameConfig struct {
	Width    int
	Height   int
	MaxTicks int
}
func InitializeGame(config GameConfig, gameUpdates chan<- GameInstance) {
    gameInstance := NewGame(int64(config.Width), int64(config.Height), config.MaxTicks)

    for gameInstance.CurrentTick < gameInstance.MaxTicks {
        gameInstance.Tick()
        gameUpdates <- *gameInstance
    }

    close(gameUpdates) // Close the channel when done to signal completion
}

type GameInstance struct {
    Id          uuid.UUID
    Width       int64
    Height      int64
    MaxTicks    int
    CurrentTick int
    Board       CompositeBoard
}
// NewGame initializes a new game with the specified configuration and generates a UUID.
func NewGame(width, height int64, maxTicks int) *GameInstance {
    return &GameInstance{
        Id:          uuid.New(), // Generate a new UUID
        Width:       width,
        Height:      height,
        MaxTicks:    maxTicks,
        CurrentTick: 0,
        Board:       NewCompositeBoard(width, height),
    }
}
// Tick progresses the game state by one tick, updating the board and publishing the state.
func (g *GameInstance) Tick() {
    // Increment the game tick    
    g.CurrentTick++

    // Logic to update the board state goes here. You might implement rules of Conway's Game of Life or other logic.
    // For demonstration, let's just toggle a specific cell on each tick.
    err := g.Board.Toggle(
        Coord{
            X: int64(g.CurrentTick % int(g.Width)), 
            Y: int64(g.CurrentTick % int(g.Height)),
        },
    )
    if err != nil {
        fmt.Println("Error toggling cell: ", err)
    }

    // Example: Check if we've reached MaxTicks and reset or end the game as needed.
    if g.CurrentTick >= g.MaxTicks {
        fmt.Println("Reached maximum number of ticks.")
        // Reset or handle the end of the game as needed.
    }
}

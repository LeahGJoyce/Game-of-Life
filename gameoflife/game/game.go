package game

import (
	"fmt"

	"github.com/google/uuid"
)


type GameConfig struct {
	Id       uuid.UUID
	Width    int
	Height   int
	MaxTicks int
}

func InitializeGame(config GameConfig, gameUpdates chan<- GameInstance) {
	gameInstance := NewGame(config.Id, int64(config.Width), int64(config.Height), config.MaxTicks)
	// Karel starting state
	// gameInstance.Board.Toggle(Coord{X: 7, Y: 4})
	// gameInstance.Board.Toggle(Coord{X: 12, Y: 4})
	// gameInstance.Board.Toggle(Coord{X: 7, Y: 5})
	// gameInstance.Board.Toggle(Coord{X: 8, Y: 5})
	// gameInstance.Board.Toggle(Coord{X: 9, Y: 5})
	// gameInstance.Board.Toggle(Coord{X: 10, Y: 5})
	// gameInstance.Board.Toggle(Coord{X: 11, Y: 5})
	// gameInstance.Board.Toggle(Coord{X: 12, Y: 5})
	// gameInstance.Board.Toggle(Coord{X: 7, Y: 6})
	// gameInstance.Board.Toggle(Coord{X: 12, Y: 6})
	// gameInstance.Board.Toggle(Coord{X: 7, Y: 10})
	// gameInstance.Board.Toggle(Coord{X: 8, Y: 10})
	// gameInstance.Board.Toggle(Coord{X: 9, Y: 10})
	// gameInstance.Board.Toggle(Coord{X: 10, Y: 10})
	// gameInstance.Board.Toggle(Coord{X: 11, Y: 10})
	// gameInstance.Board.Toggle(Coord{X: 12, Y: 10})
	// gameInstance.Board.Toggle(Coord{X: 6, Y: 11})
	// gameInstance.Board.Toggle(Coord{X: 13, Y: 11})
	// gameInstance.Board.Toggle(Coord{X: 5, Y: 12})
	// gameInstance.Board.Toggle(Coord{X: 14, Y: 12})
	// gameInstance.Board.Toggle(Coord{X: 6, Y: 13})
	// gameInstance.Board.Toggle(Coord{X: 13, Y: 13})
	// gameInstance.Board.Toggle(Coord{X: 7, Y: 14})
	// gameInstance.Board.Toggle(Coord{X: 8, Y: 14})
	// gameInstance.Board.Toggle(Coord{X: 9, Y: 14})
	// gameInstance.Board.Toggle(Coord{X: 10, Y: 14})
	// gameInstance.Board.Toggle(Coord{X: 11, Y: 14})
	// gameInstance.Board.Toggle(Coord{X: 12, Y: 14})

	// acorn starting state
	gameInstance.Board.Toggle(Coord{X: 7, Y: 4})
	gameInstance.Board.Toggle(Coord{X: 9, Y: 5})
	gameInstance.Board.Toggle(Coord{X: 6, Y: 6})
	gameInstance.Board.Toggle(Coord{X: 7, Y: 6})
	gameInstance.Board.Toggle(Coord{X: 10, Y: 6})
	gameInstance.Board.Toggle(Coord{X: 11, Y: 6})
	gameInstance.Board.Toggle(Coord{X: 12, Y: 6})

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
func NewGame(Id uuid.UUID, width, height int64, maxTicks int) *GameInstance {
	return &GameInstance{
		Id:          Id,
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

	newBoard := NewCompositeBoard(g.Width, g.Height)
	for iY := range g.Height {
		for iX := range g.Width {
			currentCoord := Coord{X: iX, Y: iY}
			surroundingCoords := getSurroundingCoordinates(g.Height, g.Width, currentCoord)
			liveNeighbors, err := countLiveNeighbors(&g.Board, surroundingCoords)
			if err != nil {
				fmt.Printf("Failed counting live neighbors at (%d, %d)", iX, iY)
			}
			// Any live cell with fewer than two live neighbors dies, as if by underpopulation.
			// Any live cell with two or three live neighbors lives on to the next generation.
			// Any live cell with more than three live neighbors dies, as if by overpopulation.
			// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.

			// Apply the rules of the game
			currentCellState, _ := g.Board.Get(currentCoord)
			if currentCellState && liveNeighbors < 2 {
				newBoard.Set(currentCoord, false) // Dies by underpopulation
			} else if currentCellState && (liveNeighbors == 2 || liveNeighbors == 3) {
				newBoard.Set(currentCoord, true) // Lives on to the next generation
			} else if currentCellState && liveNeighbors > 3 {
				newBoard.Set(currentCoord, false) // Dies by overpopulation
			} else if !currentCellState && liveNeighbors == 3 {
				newBoard.Set(currentCoord, true) // Becomes alive by reproduction
			} else {
				newBoard.Set(currentCoord, currentCellState) // No change
			}
		}
	}

	g.Board = newBoard

	if g.CurrentTick >= g.MaxTicks {
		fmt.Println("Reached maximum number of ticks.")
	}
}
func getSurroundingCoordinates(height, width int64, coord Coord) []Coord {
	var surroundingCoords []Coord

	// Define the relative positions of surrounding cells
	relativeCoords := [][]int64{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, relCoord := range relativeCoords {
		newX := coord.X + relCoord[0]
		newY := coord.Y + relCoord[1]

		// Adjust newX and newY to wrap around if out of bounds
		if newX < 0 {
			newX += width
		} else if newX >= width {
			newX -= width
		}
		if newY < 0 {
			newY += height
		} else if newY >= height {
			newY -= height
		}

		surroundingCoords = append(surroundingCoords, Coord{X: newX, Y: newY})
	}

	return surroundingCoords
}
func countLiveNeighbors(board *CompositeBoard, coords []Coord) (int, error) {
	count := 0
	for _, coord := range coords {
		// Get the value of the current coordinate
		value, err := board.Get(coord)
		if err != nil {
			return 0, err
		}
		if value {
			count++
		}
	}
	return count, nil
}

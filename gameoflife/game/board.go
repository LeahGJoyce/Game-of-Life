package game

import (
	"errors"
	"fmt"
)

// Bitboard represents a 8x8 board using a 64-bit integer where each bit
// corresponds to a square on the board.
type Bitboard uint64

const TileWidth = 8

func (b Bitboard) String() string {
	output := ""
	for y := range TileWidth {
		for x := range TileWidth {
			sym, err := b.Get(Coord{X: int64(x), Y: int64(y)})
			if err != nil {
				fmt.Printf("Bitboard.String failed when printing (%d, %d)", x, y)
				output = output + fmt.Sprintf(" # ")
			}
			if sym {
				output = output + fmt.Sprintf(" 1 ")
			} else {
				output = output + fmt.Sprintf(" 0 ")
			}
		}
		output = output + "\n"
	}
	output = output + "\n"
	return output
}

// Get returns true or false based on the value of the bitboard at location (x, y)
func (b Bitboard) Get(coord Coord) (bool, error) {
	if coord.X > TileWidth-1 || coord.Y > TileWidth-1 {
		// throw an error as x, y should be in range [0, 7] for an 8x8 bitboard
		return false, errors.New(fmt.Sprintf("Bitboard.Get received invalid coordinates (%d, %d)", coord.X, coord.Y))
	}
	pos := coord.Y*TileWidth + coord.X // calculate the bit position
	return (b & (Bitboard(1) << pos)) != 0, nil
}

func (b *Bitboard) Toggle(coord Coord) error {
	if coord.X > TileWidth-1 || coord.Y > TileWidth-1 {
		// throw an error as x, y should be in range [0, 7] for an 8x8 bitboard
		return errors.New(fmt.Sprintf("Bitboard.ToggleSquare recieved invalid coordinates (%d, %d)", coord.X, coord.Y))
	}
	pos := coord.Y*TileWidth + coord.X // calculate the bit position
	*b ^= Bitboard(1) << pos
	return nil
}

// CompositeBoard represents a larger board composed of multiple 8x8 Bitboard tiles.
type CompositeBoard struct {
	BoardWidth  int64
	BoardHeight int64
	Tiles       [][]Bitboard
}

func (b CompositeBoard) String() string {
	output := ""
	for iY := range b.BoardHeight {
		for iX := range b.BoardWidth {
			sym, err := b.Get(Coord{X: iX, Y: iY})
			if err != nil {
				fmt.Printf("Bitboard.String failed when printing (%d, %d)", iX, iY)
				output = output + fmt.Sprintf(" # ")
			}
			if sym {
				output = output + fmt.Sprintf(" 1 ")
			} else {
				output = output + fmt.Sprintf(" 0 ")
			}

		}
		output = output + "\n"
	}
	output = output + "\n"
	return output
}

func (b *CompositeBoard) Get(coord Coord) (bool, error) {
	if coord.X >= b.BoardWidth || coord.Y >= b.BoardHeight {
		return false, errors.New(fmt.Sprintf("CompositeBoard.Get received invalid coordinates (board: %d, %d)", coord.X, coord.Y))
	}
	boardCoord, bitCoord := convertCoordinates(coord)

	return b.Tiles[boardCoord.X][boardCoord.Y].Get(bitCoord)
}

func (b *CompositeBoard) Toggle(coord Coord) error {
	if coord.X >= b.BoardWidth || coord.Y >= b.BoardHeight {
		return errors.New(fmt.Sprintf("CompositeBoard.Toggle received invalid coordinates (board: %d, %d)", coord.X, coord.Y))
	}
	boardCoord, bitCoord := convertCoordinates(coord)

	return b.Tiles[boardCoord.X][boardCoord.Y].Toggle(bitCoord)
}

// NewCompositeBoard creates and returns a new CompositeBoard with the specified dimensions.
func NewCompositeBoard(width, height int64) CompositeBoard {
	// Divide width and height by TileWidth to get the correct amount of BitBoards
	boardWidth := width / TileWidth
	boardHeight := height / TileWidth

	// Check for any remaining cells and add an extra BitBoard if necessary
	if width%TileWidth != 0 {
		boardWidth++
	}
	if height%TileWidth != 0 {
		boardHeight++
	}

	tiles := make([][]Bitboard, boardHeight)
	for i := range tiles {
		tiles[i] = make([]Bitboard, boardWidth)
		// Initialize each tile as a Bitboard with all bits set to 0
		for j := range tiles[i] {
			tiles[i][j] = Bitboard(0)
		}
	}
	return CompositeBoard{
		BoardWidth:  width,
		BoardHeight: height,
		Tiles:       tiles,
	}
}

type Coord struct {
	X int64
	Y int64
}

func convertCoordinates(coord Coord) (Coord, Coord) {
	boardX := coord.X / TileWidth
	boardY := coord.Y / TileWidth

	bitX := coord.X % TileWidth
	bitY := coord.Y % TileWidth

	return Coord{X: boardX, Y: boardY}, Coord{X: bitX, Y: bitY}
}

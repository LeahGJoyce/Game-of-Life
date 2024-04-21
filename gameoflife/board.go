package main

import (
	"gameoflife/proto"
	v4 "github.com/google/uuid"
)

func RenderBoard(board *proto.Board) string {
	// TODO: Update this to return a string representation of the board
	return "<THE BOARD>"
}

func CreateBlankBoard(boardHeight int, boardWidth int) *proto.Board {
	board := proto.Board{
		Id:     v4.New().String(),
		Height: int64(boardHeight),
		Width:  int64(boardWidth),
		Rows:   make([]*proto.Board_Row, boardHeight),
	}

	for i := range board.Rows {
		row := proto.Board_Row{Columns: make([]bool, boardWidth)}
		for j := range row.Columns {
			row.Columns[j] = false
		}
		board.Rows[i] = &row
	}
	return &board
}

// CreateExampleBoard function is unused. It demonstrates how to
// create a board with a predefined state
func CreateExampleBoard() proto.Board {
	row1 := proto.Board_Row{Columns: []bool{true, true, true}}
	row2 := proto.Board_Row{Columns: []bool{true, false, true}}
	row3 := proto.Board_Row{Columns: []bool{true, true, true}}
	return proto.Board{Id: v4.New().String(), Height: 3, Width: 3, Rows: []*proto.Board_Row{&row1, &row2, &row3}}
}

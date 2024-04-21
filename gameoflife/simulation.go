package main

import (
	"gameoflife/proto"
	v4 "github.com/google/uuid"
	"time"
)

func CreateSimulation(boardHeight int, boardWidth int, maxTicks int, msPerTick int) *proto.Simulation {
	board := proto.Board{
		Id:     v4.New().String(),
		Height: int64(boardHeight),
		Width:  int64(boardWidth),
		Column: make([]*proto.Board_Row, boardHeight),
	}

	for i := range board.Column {
		row := proto.Board_Row{Cell: make([]bool, boardWidth)}
		for j := range row.Cell {
			row.Cell[j] = false
		}
		board.Column[i] = &row
	}

	simulation := proto.Simulation{
		Id:          v4.New().String(),
		MaxTicks:    int64(maxTicks),
		MsPerTick:   int64(msPerTick),
		CurrentTick: 0,
		Board:       &board,
	}
	return &simulation
}

func RunSimulation(simulation *proto.Simulation) {
	for i := range simulation.MaxTicks {
		simulation.CurrentTick = i
		// TODO: Add logic for updating the simulation and board
		time.Sleep(time.Duration(simulation.MsPerTick) * time.Millisecond)
	}
}

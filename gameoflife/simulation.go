package main

import (
	"gameoflife/proto"
	v4 "github.com/google/uuid"
	"time"
)

func CreateSimulation(boardHeight int, boardWidth int, maxTicks int, msPerTick int) *proto.Simulation {
	simulation := proto.Simulation{
		Id:          v4.New().String(),
		MaxTicks:    int64(maxTicks),
		MsPerTick:   int64(msPerTick),
		CurrentTick: 0,
		Board:       CreateBlankBoard(boardHeight, boardWidth),
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

package main

import (
	"fmt"
	"gameoflife/proto"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	simulation *proto.Simulation
	menu       []string
	cursor     int
}

func initialModel() model {
	menu := GetMenu()
	cursor := 0
	var simulation proto.Simulation
	return model{&simulation, menu, cursor}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

type TickMsg time.Time

// Send a message every second.
func tickEvery(msPerTick int64) tea.Cmd {
	return tea.Every(time.Duration(msPerTick)*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.menu)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) select an option
		case "enter", " ":
			if m.menu[m.cursor] == "Start Simulation" {
				BoardHeight := 10
				BoardWidth := 10
				MaxTicks := 25
				MSPerTick := 250
				m.simulation = CreateSimulation(BoardHeight, BoardWidth, MaxTicks, MSPerTick)
				go RunSimulation(m.simulation)
				return m, tickEvery(m.simulation.MsPerTick)
			}
			if m.menu[m.cursor] == "Quit" {
				return m, tea.Quit
			}
		}
	case TickMsg:
		if m.simulation.CurrentTick == m.simulation.MaxTicks-1 {
			simulation := proto.Simulation{}
			m.simulation = &simulation
		}
		return m, tickEvery(m.simulation.MsPerTick)
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Conway's Game of Life\n\n"

	if m.simulation.Id == "" {
		s += RenderMenu(m)
	} else {
		s += fmt.Sprintf("Tick: %d\n", m.simulation.CurrentTick+1)
		s += RenderBoard(m.simulation.Board)
	}

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"gameoflife/game"
	"gameoflife/publisher"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func startHandler(w http.ResponseWriter, r *http.Request) {
	// Configuration for the game, potentially extracted from the request
	config := game.GameConfig{
		Width:    10, // Example values, adjust as needed
		Height:   10,
		MaxTicks: 20,
	}

	// Create a channel for game state updates
	gameUpdates := make(chan game.GameInstance)

	// Start the game simulation in a separate goroutine
	go game.InitializeGame(config, gameUpdates)

	// Listen for game state updates and publish them
	go func() {
		for gameInstance := range gameUpdates {
			// Convert UUID to string
			gameIDStr := gameInstance.Id.String()
			gameBoardStr := gameInstance.Board.String()

			// Now use gameIDStr as part of the stream name or in messages
			fmt.Printf("Publishing update for game %s: %s\n", gameIDStr, gameBoardStr)
			publisher.PublishMessage(gameIDStr, gameBoardStr)
		}
	}()

	// Respond to the HTTP request immediately
	fmt.Fprintf(w, "Game started successfully!")
}

func main() {
	// Serve static files from the web directory
	fs := http.FileServer(http.Dir("web/"))
	http.Handle("/", fs)

	// Setup signal handling for a graceful shutdown.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	// Set up other routes, e.g., for your API endpoints
	http.HandleFunc("/start", startHandler)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	// Wait for interrupt signal.
	<-stopChan

	// Begin cleanup sequence.
	publisher.CleanupProducers()
	publisher.CleanupEnvironment()

	fmt.Println("Application shutdown gracefully")
}

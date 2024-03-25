package main

import (
	"fmt"
	"gameoflife/game"
	"gameoflife/publisher"
	"net/http"
)

func startHandler(w http.ResponseWriter, r *http.Request) {
	// Configuration for the game, potentially extracted from the request
	config := game.GameConfig{
		Width:    20, // Example values, adjust as needed
		Height:   20,
		MaxTicks: 30,
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
			fmt.Printf("Publishing update for game %s\n%s\n", gameIDStr, gameBoardStr)
			env := publisher.GetEnv()
			publisher.PublishMessage(env, gameIDStr, gameBoardStr)
		}
	}()

	// Respond to the HTTP request immediately
	fmt.Fprintf(w, "Game started successfully!")
}

func main() {
	// Serve static files from the web directory
	fs := http.FileServer(http.Dir("web/"))
	http.Handle("/", fs)

	// Set up other routes, e.g., for your API endpoints
	http.HandleFunc("/start", startHandler)

	fmt.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	// Cleanup Publisher
	env := publisher.GetEnv()
	err := env.Close()
	publisher.CheckErr(err)
}

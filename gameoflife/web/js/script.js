document.getElementById('startSimulationBtn').addEventListener('click', function() {
    startSimulation();
});

function startSimulation() {
    // Define the endpoint URL
    const endpoint = 'http://localhost:8080/start'; // Update with the correct URL/port if needed

    // Define the simulation start parameters here. Example:
    const simulationParams = {
        initialState: "...", // Your initial state data here
        boardSize: "...", // Your board size here
        tickCount: "..." // Your tick count here
    };

    // Make the API call to start the simulation
    fetch(endpoint, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(simulationParams),
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        }
        throw new Error('Network response was not ok.');
    })
    .then(data => console.log(data))
    .catch(error => console.error('There has been a problem with your fetch operation:', error));
}

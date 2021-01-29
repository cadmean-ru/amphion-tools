window.amphionDebugger = {
    start: (port) => {
        let socket = new WebSocket(`ws://localhost:${port}/ws`);

        console.log("Attempting to connect to debug server...");

        socket.onopen = () => {
            console.log("Connected to debug server.");
            socket.send("hello");
        };

        socket.onmessage = e => {
            console.log(`Data from debug: ${e.data}`);

            if (e.data === "refresh") {
                socket.close(1000, "work done: page reload")
                document.location.reload();
            }
        };

        socket.onclose = event => {
            console.log("Socket Closed Connection: ", event);
            socket.send("close")
        };

        socket.onerror = error => {
            console.log("Socket Error: ", error);
        };
    }
}
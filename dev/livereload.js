const eventSource = new EventSource("http://localhost:3001/livereload")

eventSource.onmessage = (event) => {
    console.log("event:", event)
    if (event.data == "fsChange") {
        window.location.reload()
    }
}
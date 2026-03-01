const eventSource = new EventSource("/dev/livereload")

eventSource.onmessage = (event) => {
    if (event.data == "fsChange") {
        window.location.reload()
    }
}
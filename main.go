package main

import (
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/event-stream", eventHandler)
	http.ListenAndServe(":8080", nil)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	data := []string{
		"hey",
		"my",
		"name",
		"is",
		"shivam",
		"singh",
	}

	for _, d := range data {

		w.Write([]byte(d))
		w.(http.Flusher).Flush()
		time.Sleep(time.Second * 3)
	}
}

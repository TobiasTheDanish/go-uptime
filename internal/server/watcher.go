package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Starts go-routine that pings the endpoints passed in, and watches for unresponsive services.
// Parameter endpoints is a slice of strings containing the endpoints for the status checks.
func StartWatch(endpoints []string) {
	if len(endpoints) == 0 {
		fmt.Fprintln(os.Stderr, "No endpoints provided to status watcher. Returning")
		return
	}

	go watchEndpoints(os.Stdout, os.Stderr, endpoints)
}

func watchEndpoints(out, errOut io.Writer, endpoints []string) {
	t := time.Now().Unix()

	for {
		delta := time.Duration(time.Now().Unix() - t)

		if delta >= 5*time.Second {
			for _, url := range endpoints {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					fmt.Fprintf(errOut, "Error creating new request: %s\n", err)
					return
				}

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprintf(errOut, "Error executing request to %s: %s\n", url, err)
					return
				}
				defer res.Body.Close()

				if res.StatusCode != 200 {
					fmt.Fprintf(errOut, "Server at %s is down\n", url)
				} else {
					fmt.Fprintf(out, "Server at %s is up\n", url)
				}
			}
			t = time.Now().Unix()
		}
	}
}

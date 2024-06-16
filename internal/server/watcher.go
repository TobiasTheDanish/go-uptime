package server

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Status struct {
	Time        int64
	IsUp        bool
	ReponseTime int64
}

var StatusMap map[string][]Status

// Starts go-routine that pings the endpoints passed in, and watches for unresponsive services.
// Parameter endpoints is a slice of strings containing the endpoints for the status checks.
func StartWatch(endpoints []string) {
	if len(endpoints) == 0 {
		StatusMap = nil
		fmt.Fprintln(os.Stderr, "No endpoints provided to status watcher. Returning")
		return
	}

	StatusMap = make(map[string][]Status, len(endpoints))
	for _, ep := range endpoints {
		StatusMap[ep] = make([]Status, 0, 0)
	}

	go watchEndpoints(endpoints)
}

func watchEndpoints(endpoints []string) {
	t := time.Now().Unix()

	for {
		delta := time.Duration(time.Now().Unix() - t)

		if delta >= 5 {
			for _, url := range endpoints {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error creating new request: %s\n", err)
					return
				}

				status := Status{
					Time: time.Now().UnixMilli(),
				}
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error executing request to %s: %s\n", url, err)
					return
				}
				defer res.Body.Close()

				status.ReponseTime = time.Now().UnixMilli() - status.Time
				status.IsUp = res.StatusCode == 200
				StatusMap[url] = append(StatusMap[url], status)
			}
			t = time.Now().Unix()
		}
	}
}

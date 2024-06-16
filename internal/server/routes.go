package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.GET("/", s.GetStatusHandler)

	return e
}

type State struct {
	State                 string
	StartTimeMS           int64
	DurationMS            int64
	AverageResponseTimeMs int64
}

func (s *Server) GetStatusHandler(c echo.Context) error {
	stateMap := make(map[string][]State)

	for url, statusSlice := range StatusMap {
		stateMap[url] = make([]State, 0)
		if len(statusSlice) == 0 {
			break
		}

		var prevUp bool
		var currentState State
		count := 1
		responseTimes := int64(statusSlice[0].ReponseTime)

		for _, status := range statusSlice {
			if currentState == (State{}) || status.IsUp != prevUp {
				if currentState != (State{}) {
					currentState.DurationMS = status.Time - currentState.StartTimeMS
					currentState.AverageResponseTimeMs = responseTimes / int64(count)
					stateMap[url] = append(stateMap[url], currentState)
					count = 1
					responseTimes = status.ReponseTime
				}

				currentState = State{
					State:       isUpString(status.IsUp),
					StartTimeMS: status.Time,
				}
			} else {
				responseTimes += status.ReponseTime
				count += 1
			}
			prevUp = status.IsUp
		}
		currentState.DurationMS = statusSlice[len(statusSlice)-1].Time - currentState.StartTimeMS
		currentState.AverageResponseTimeMs = responseTimes / int64(count)
		stateMap[url] = append(stateMap[url], currentState)
	}

	return c.JSON(http.StatusOK, stateMap)
}

func isUpString(isUp bool) string {
	if isUp {
		return "Up"
	}
	return "Down"
}

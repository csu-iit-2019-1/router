package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	eventHost            = ""
	eventPort            = ""
	getEventPathTemplate = "/events/%d"
)

type Event struct {
	EventId     int       `json:"eventId"`
	Name        string    `json:"name"`
	City        int       `json:"city"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	FreeSpace   int       `json:"free_space"`
}

func GetEventsByIds(ids []int) map[int]*Event {
	mtx := new(sync.Mutex)
	events := make(map[int]*Event)

	wg := new(sync.WaitGroup)

	for _, eventId := range ids {
		wg.Add(1)
		go func(id int, group *sync.WaitGroup) {
			defer group.Done()

			req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s%s", eventHost, eventPort, fmt.Sprintf(getEventPathTemplate, id)), nil)
			if err != nil {
				panic(err)
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			event := new(Event)
			err = json.Unmarshal(body, event)
			if err != nil {
				panic(err)
			}
			mtx.Lock()
			events[id] = event
			mtx.Unlock()

		}(eventId, wg)
	}
	return events
}

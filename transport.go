package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	transportHost    = ""
	transportPort    = ""
	getTransportPath = ""
)

type TKey struct {
	time.Time
	startPoint int
	endPoint   int
}

type Transport struct {
	TransportId   int     `json:"transport_id"`
	TransportType string  `json:"transport_type"` //[ "aircraft", "train", "bus" ]
	Name          string  `json:"name"`
	StartPoint    int     `json:"start_point"`
	EndPoint      int     `json:"end_point"`
	DepartureTime int     `json:"departure_time"`
	ArriveTime    int     `json:"arrive_time"`
	Price         float64 `json:"price"`
}

func GetAllTransport(requirements *RouteRequirements) map[TKey]*Transport {
	mtx := new(sync.Mutex)
	transport := make(map[TKey]*Transport)
	wg := new(sync.WaitGroup)

	for timeKey := requirements.StartDate; timeKey.Before(requirements.EndDate); timeKey.AddDate(0, 0, 1) {
		citiesIds := []int{requirements.DepartureCityId}
		for _, cityReq := range requirements.CitiesRequirements {
			citiesIds = append(citiesIds, cityReq.CityId)
		}
		for _, startPoint := range citiesIds {
			wg.Add(1)

			go func(tKey time.Time, startP int, endPs []int, group *sync.WaitGroup) {
				defer wg.Done()
				reqData := struct {
					DepartureDate int      `json:"departureDate"`
					TransportType []string `json:"transportType"`
					StartPoint    int      `json:"startPoint"`
					EndPoints     []int    `json:"endPoints"`
				}{
					DepartureDate: int(tKey.Unix()),
					TransportType: requirements.TransportTypes,
					StartPoint:    startP,
					EndPoints:     endPs,
				}
				jsonStr, err := json.Marshal(reqData)
				if err != nil {
					panic(err)
				}
				req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s%s", transportHost, transportPort, getTransportPath), bytes.NewBuffer(jsonStr))
				if err != nil {
					panic(err)
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				fmt.Println("response Status:", resp.Status)
				fmt.Println("response Headers:", resp.Header)
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("response Body:", string(body))

				respTransport := []*Transport{}
				err = json.Unmarshal(body, &respTransport)
				if err != nil {
					panic(err)
				}
				for _, t := range respTransport {
					if t.TransportId != -1 {
						mtx.Lock()
						transport[TKey{tKey, startP, t.EndPoint}] = t
						mtx.Unlock()
					}
				}

			}(timeKey, startPoint, citiesIds, wg)
		}
	}
	return transport
}

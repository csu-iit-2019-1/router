package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	hotelHost            = ""
	hotelPort            = ""
	getCheapestHotelPath = ""
)

type Hotel struct {
	Id           int     `json:"id"`
	CityId       int     `json:"city"`
	Name         string  `json:"name"`
	MainPhotoUrl string  `json:"mainPhotoUrl"`
	Price        float64 `json:"price"`
	SeaIsNear    bool    `json:"seaIsNear"`
	Breakfast    bool    `json:"breakfast"`
	Stars        float64 `json:"stsars"`
}

func GetCheapestHotels(requirements *RouteRequirements) map[int]*Hotel {

	mtx := new(sync.Mutex)
	hotels := make(map[int]*Hotel)

	wg := new(sync.WaitGroup)

	for _, city := range requirements.CitiesRequirements {
		wg.Add(1)
		go func(cityReq *CityRequirements, group *sync.WaitGroup) {
			defer group.Done()
			reqData := struct {
				Date        string `json:"date"`
				CityId      int    `json:"cityId"`
				Stars       int    `json:"stars"`
				Breakfast   bool   `json:"breakfast"`
				SeaIsNearby bool   `json:"seaIsNearby"`
			}{
				Date:        "",
				CityId:      cityReq.CityId,
				Stars:       requirements.MinHotelStars,
				Breakfast:   cityReq.HotelBreakfast,
				SeaIsNearby: cityReq.HotelBreakfast,
			}
			jsonStr, err := json.Marshal(reqData)
			if err != nil {
				panic(err)
			}
			req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%s%s", hotelHost, hotelPort, getCheapestHotelPath), bytes.NewBuffer(jsonStr))
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
			hotelInf := struct {
				ShortInfAboutHotel []*Hotel `json:"shortInfAboutHotel"`
			}{}
			err = json.Unmarshal(body, &hotelInf)
			if err != nil {
				panic(err)
			}
			mtx.Lock()
			hotels[hotelInf.ShortInfAboutHotel[0].CityId] = hotelInf.ShortInfAboutHotel[0]
			mtx.Unlock()
		}(city, wg)
	}
	wg.Wait()
	return hotels
}

package main

import (
	"time"
)

type CityRequirements struct {
	CityId         int         `json:"city_id"`
	TotalDaysIn    int         `json:"total_days_in"`
	DatesIn        []time.Time `json:"dates_in"`
	EventIds       []int       `json:"event_ids"`
	HotelBreakfast bool        `json:"hotel_breakfast"`
	HotelSeaNearby bool        `json:"hotel_sea_nearby"`
}

type RouteRequirements struct {
	StartDate          time.Time           `json:"start_date"`
	EndDate            time.Time           `json:"end_date"`
	PersonAmount       int                 `json:"person_amount"`
	DepartureCityId    int                 `json:"departure_city_id"`
	TransportTypes     []string            `json:"transport_types"`
	MinHotelStars      int                 `json:"min_hotel_stars"`
	CitiesRequirements []*CityRequirements `json:"cities_requirements"`
}

type RoutePoint struct {
	CityId    int        `json:"city_id"`
	Hotel     *Hotel     `json:"hotel"`
	Transport *Transport `json:"transport"`
	Events    []*Event   `json:"events"`
}

type Route struct {
	Points           []*RoutePoint `json:"points"`
	ReverseTransport *Transport    `json:"reverse_transport"`
	TotalPrice       float64       `json:"total_price"`
}

func buildRoute(requirements *RouteRequirements) *Route {
	//route := new(Route)
	//
	//cityHotels := GetCheapestHotels(requirements)
	//transports := GetAllTransport(requirements)
	//eventsIds := []int{}
	//for _, cityR := range requirements.CitiesRequirements {
	//	for _, id := range cityR.EventIds {
	//		eventsIds = append(eventsIds, id)
	//	}
	//}
	//events := GetEventsByIds(eventsIds)
	//
	//
	//for _, cityR := range requirements.CitiesRequirements {
	//	cityHotel := cityHotels[cityR.CityId]
	//	cityT
	//}
	return &Route{}
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"googlemaps.github.io/maps"
)

type Body struct {
	Origin      string `json:"origen,omitempty"`
	Destination string `json:"destino,omitempty"`
}

func GetDirectionsEndpoint(w http.ResponseWriter, req *http.Request) {
	var body Body
	_ = json.NewDecoder(req.Body).Decode(&body)

	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyCLJS1j7S7bTFHAk8oiplSt748Ivb2AGw4"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	r := &maps.DirectionsRequest{
		Origin:      body.Origin,
		Destination: body.Destination,
	}
	resp, _, err := c.Directions(context.Background(), r)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	data := []byte(`{"routes:":[`)

	json.NewDecoder(req.Body).Decode(&resp)

	for x := 0; x < len(resp[0].Legs[0].Steps); x++ {
		data = append(data, "{\"lat\":"...)
		data = append(data, strconv.FormatFloat(resp[0].Legs[0].Steps[x].StartLocation.Lat, 'f', 5, 64)...)
		data = append(data, ", "...)
		data = append(data, "\"lon\":"...)
		data = append(data, strconv.FormatFloat(resp[0].Legs[0].Steps[x].StartLocation.Lng, 'f', 5, 64)...)

		if x == len(resp[0].Legs[0].Steps)-1 {
			data = append(data, "} "...)
		} else {
			data = append(data, "}, "...)
		}
	}

	data = append(data, "]}"...)

	fmt.Fprintf(w, string(data))

}

func GetRestaurantsEndpoint(w http.ResponseWriter, req *http.Request) {
	var body Body
	_ = json.NewDecoder(req.Body).Decode(&body)

	client, err := maps.NewClient(maps.WithAPIKey("AIzaSyBmelZAhVTODrw_gjtueTuHEs9Aka_z9nM"))
	if err != nil {
		log.Fatalf("Fatal Error: %s", err)
	}

	body_detail := &maps.GeocodingRequest{
		Address: body.Origin,
	}

	body_response, _ := client.Geocode(context.Background(), body_detail)

	r := &maps.NearbySearchRequest{

		Location: &body_response[0].Geometry.Location,
		Radius:   100,
		Type:     maps.PlaceTypeRestaurant,
	}

	restaurants, _ := client.NearbySearch(context.Background(), r)
	json.NewDecoder(req.Body).Decode(&restaurants)

	data := []byte(`{"restaurantes:":[`)

	for x := 0; x < len(restaurants.Results); x++ {
		data = append(data, "{\"nombre\":\""...)
		data = append(data, restaurants.Results[x].Name...)
		data = append(data, "\", "...)
		data = append(data, "\"lat\":"...)
		data = append(data, strconv.FormatFloat(restaurants.Results[x].Geometry.Location.Lat, 'f', 5, 64)...)
		data = append(data, ", "...)
		data = append(data, "\"lon\":"...)
		data = append(data, strconv.FormatFloat(restaurants.Results[x].Geometry.Location.Lng, 'f', 5, 64)...)
		if x == len(restaurants.Results)-1 {
			data = append(data, "} "...)
		} else {
			data = append(data, "}, "...)
		}
	}

	data = append(data, "]}"...)
	fmt.Fprintf(w, string(data))

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ejercicio1", GetDirectionsEndpoint).Methods("POST")
	router.HandleFunc("/ejercicio2", GetRestaurantsEndpoint).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

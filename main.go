package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rwcarlsen/goexif/exif"
)

func main() {
    lambda.Start(GoogleMapApiExample)
}

type Event struct {
    Image string `json:"image"`
}

func GoogleMapApiExample(e Event) (places []PlaceType, err error) {
    b, err := base64.StdEncoding.DecodeString(e.Image)
    if err != nil {
        return
    }
    r := bytes.NewReader(b)
    lat, long, err := getLocation(r)
    if err != nil {
        return
    }

    p, err := fetchPlaces(lat, long)
    if err != nil {
        return
    }

    for _, place := range p.Results {
        for _, t := range place.Types {
            if t == "cafe" || t == "restaurant" {
                places = append(places, place)
            }
        }
    }

    return
}

type PlaceType struct {
    Name string `json:"name"`
    PlaceId string `json:"place_id"`
    Types []string `json:"types"`
}

type PlacesType struct {
    Results []PlaceType `json:"results"`
}

func fetchPlaces(lat float64, long float64) (PlacesType, error) {
    url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%g,%g&language=ja&radius=10&key=%s", lat, long, os.Getenv("API_KEY"))
    resp, err := http.Get(url)
    if err != nil {
        return PlacesType{}, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return PlacesType{}, err
    }

    var places PlacesType
    err = json.Unmarshal(body, &places)

    return places, err
}

func getLocation(r *bytes.Reader) (lat float64, long float64, err error) {
    x, err := exif.Decode(r)
    if err != nil {
        return
    }
    lat, long, err = x.LatLong()
    return
}

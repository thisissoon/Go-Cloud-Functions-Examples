package location

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/genproto/googleapis/type/latlng"
)

type LatLong struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Response struct {
	Status int     `json:"status"`
	Result LatLong `json:"result"`
}

func GetLatLong(postcode string) (*latlng.LatLng, error) {
	apiUrl := "https://api.postcodes.io/postcodes/"
	resp, err := http.Get(apiUrl + postcode)
	if err != nil {
		return &latlng.LatLng{}, fmt.Errorf("error getting postcode: %v", err)
	}
	if resp.StatusCode != 200 {
		return &latlng.LatLng{}, fmt.Errorf("error making postcode request. Status is :%v", resp.StatusCode)
	}
	var res Response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &latlng.LatLng{}, fmt.Errorf("error reading json bytes: %v", err)
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &res)
	if err != nil {
		return &latlng.LatLng{}, fmt.Errorf("error unmarshaling json body: %v", err)
	}
	ll := &latlng.LatLng{
		Latitude:  res.Result.Latitude,
		Longitude: res.Result.Longitude,
	}
	return ll, nil
}

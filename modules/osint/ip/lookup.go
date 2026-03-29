package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yash3605/spectre/internal/models"
)

type response struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Latitude    float32 `json:"lat"`
	Longitude   float32 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	ORG         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
}

var ipFieldOrder = []string{
	"query",
	"status",
	"country",
	"countryCode",
	"regionName",
	"region",
	"city",
	"zip",
	"latitude",
	"longitude",
	"timezone",
	"isp",
	"org",
	"as",
}

func Lookup(address string) models.Result {
	res, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", address))
	if err != nil {
		fmt.Printf("There's an error looking up for the IP: %v", err)
		return models.Result{
			Title:  address,
			Data:   map[string]string{},
			Status: models.StateError,
		}
	}

	defer res.Body.Close()

	resSlice, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response")
		return models.Result{
			Title:  address,
			Data:   map[string]string{},
			Status: models.StateError,
		}
	}

	var ipRes response
	err = json.Unmarshal(resSlice, &ipRes)
	if err != nil {
		fmt.Println("Error parsing response into JSON")
		return models.Result{
			Title:  address,
			Data:   map[string]string{},
			Status: models.StateError,
		}
	}

	return models.Result{
		Title: address,
		Data: map[string]string{
			"status":      ipRes.Status,
			"country":     ipRes.Country,
			"countryCode": ipRes.CountryCode,
			"region":      ipRes.Region,
			"regionName":  ipRes.RegionName,
			"city":        ipRes.City,
			"zip":         ipRes.Zip,
			"latitude":    fmt.Sprintf("%f", ipRes.Latitude),
			"longitude":   fmt.Sprintf("%f", ipRes.Longitude),
			"timezone":    ipRes.Timezone,
			"isp":         ipRes.ISP,
			"org":         ipRes.ORG,
			"as":          ipRes.AS,
			"query":       ipRes.Query,
		},
		Order:  ipFieldOrder,
		Status: models.StateSuccess,
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/alecthomas/kong"
)

var CLI struct {
	ApiKey string `help:"The API key to use with the service"`
	Lat    string `help:"Latitude of reading"`
	Lon    string `help:"Longitude of reading"`
}

type Weather struct {
	Base   string `json:"base"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Cod   int `json:"cod"`
	Coord struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
	Dt   int `json:"dt"`
	Id   int `json:"id"`
	Main struct {
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
		Temp      float64 `json:"temp"`
		TempMax   float64 `json:"temp_max"`
		TempMin   float64 `json:"temp_min"`
	} `json:"main"`
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
		Id      int    `json:"id"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
		Type    int    `json:"type"`
	} `json:"sys"`
	Timezone   int `json:"timezone"`
	Visibility int `json:"visibility"`
	Weather    []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Id          int    `json:"id"`
		Main        string `json:"main"`
	} `json:"weather"`
	Wind struct {
		Deg   int     `json:"deg"`
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Rain *struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
	Snow *struct {
		OneH float64 `json:"1h"`
	} `json:"snow"`
}

const wt = `
# TYPE weather_info gauge
# HELP weather_info Information about the location; the value equals the id of the location
# TYPE weather_return_code gauge
# HELP weather_return_code Internal API return code, presumably HTTP codes
# TYPE weather_temperature_kelvin gauge
# HELP weather_temperature_kelvin Temperature in Kelvin
# TYPE weather_temperature_feels_like_kelvin gauge
# HELP weather_temperature_feels_like_kelvin Temperature in Kelvin with human perception of temperature
# TYPE weather_wind_meters_per_second gauge
# HELP weather_wind_meters_per_second Speed in m/s
# TYPE weather_wind_direction gauge
# HELP weather_wind_direction Direction in meteorological degress
# TYPE weather_rain_1h gauge
# HELP weather_rain_1h One hour rain total
# TYPE weather_snow_1h gauge
# HELP weather_snow_1h One hour rain total
# TYPE weather_id gauge
# HELP weather_id Weather condition id
# TYPE weather_humidity_percent gauge
# HELP weather_humidity_percent Relative humidity in percent
# TYPE weather_clouds_percent gauge
# HELP weather_clouds_percent Cloud cover in percent
# TYPE weather_sun_epoch gauge
# HELP weather_sun_epoch
# TYPE weather_pressure_hectopascal gauge
# HELP weather_pressure_hectopascal
weather_info{location="{{.Name}}",lat="{{.Coord.Lat}}",lon="{{.Coord.Lon}}"} {{.Id}}
weather_return_code{id="{{.Id}}"} {{.Cod}}
weather_temperature_kelvin{id="{{.Id}}"} {{.Main.Temp}}
weather_temperature_feels_like_kelvin{id="{{.Id}}"} {{.Main.FeelsLike}}
weather_wind_meters_per_second{id="{{.Id}}"} {{.Wind.Speed}}
weather_wind_direction{id="{{.Id}}"}  {{.Wind.Deg}}
weather_rain_1h{id="{{.Id}}"} {{if .Rain}}{{.Rain.OneH}}{{else}}0{{end}}
weather_snow_1h{id="{{.Id}}"} {{if .Snow}}{{.Snow.OneH}}{{else}}0{{end}}
weather_humidity_percent{id="{{.Id}}"} {{.Main.Humidity}}
weather_clouds_percent{id="{{.Id}}"} {{.Clouds.All}}
weather_sun_epoch{id="{{.Id}}",change="sunrise"} {{.Sys.Sunrise}}
weather_sun_epoch{id="{{.Id}}",change="sunset"} {{.Sys.Sunset}}
weather_pressure_hectopascal{id="{{.Id}}"} {{.Main.Pressure}}
`

func main() {
	kong.Parse(&CLI)

	resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", CLI.Lat, CLI.Lon, CLI.ApiKey))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	w := Weather{}
	json.NewDecoder(resp.Body).Decode(&w)

	t, err := template.New("weather-prometheus").Parse(wt)
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, w)
	if err != nil {
		panic(err)
	}
}
